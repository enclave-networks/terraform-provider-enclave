package enclave

import (
	"context"
	"fmt"

	enclaveDns "github.com/enclave-networks/go-enclaveapi/data/dns"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type dnsRecordResourceType struct{}

func (d dnsRecordResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.Int64Type,
				Computed: true,
			},
			"zone_id": {
				Type:     types.Int64Type,
				Optional: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"tags": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"systems": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"notes": {
				Type:     types.StringType,
				Optional: true,
			},
			"fqdn": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

// New resource instance
func (d dnsRecordResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return dnsRecord{
		provider: *(p.(*provider)),
	}, nil
}

type dnsRecord struct {
	provider provider
}

// Create implements tfsdk.Resource
func (d dnsRecord) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !d.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, "+
				"likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan DnsRecordState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var zoneId enclaveDns.DnsZoneId
	if !plan.ZoneId.Null {
		zoneId = enclaveDns.DnsZoneId(plan.ZoneId.Value)
	} else {
		zoneId = enclaveDns.DnsZoneId(1)
	}

	dnsRecordCreate := enclaveDns.DnsRecordCreate{
		Name:    plan.Name.Value,
		ZoneId:  zoneId,
		Tags:    plan.Tags,
		Systems: plan.Systems,
		Notes:   plan.Notes.Value,
	}

	// create request
	dnsRecord, err := d.provider.client.Dns.CreateRecord(dnsRecordCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating dnsRecord in enclave",
			err.Error(),
		)
		return
	}

	setDnsRecordState(dnsRecord, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements tfsdk.Resource
func (d dnsRecord) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// read state
	var state DnsRecordState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dnsRecordId := enclaveDns.DnsRecordId(state.Id.Value)

	//call api to delete
	_, err := d.provider.client.Dns.DeleteRecord(dnsRecordId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Dns Record",
			"Could not read Id "+fmt.Sprint(dnsRecordId)+": "+err.Error(),
		)
		return
	}

	// remove resource
	resp.State.RemoveResource(ctx)
}

// Read implements tfsdk.Resource
func (d dnsRecord) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state DnsRecordState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dnsRecordId := enclaveDns.DnsRecordId(state.Id.Value)

	dnsRecord, err := d.provider.client.Dns.GetRecord(dnsRecordId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Dns Record",
			"Could not read Id "+fmt.Sprint(dnsRecordId)+": "+err.Error(),
		)
		return
	}

	setDnsRecordState(dnsRecord, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update implements tfsdk.Resource
func (d dnsRecord) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// read state to get Id
	var state DnsRecordState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get executing plan updates
	var plan DnsRecordState
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	dnsRecordId := enclaveDns.DnsRecordId(state.Id.Value)

	updateDnsRecord, err := d.provider.client.Dns.UpdateRecord(dnsRecordId, enclaveDns.DnsRecordPatch{
		Name:    plan.Name.Value,
		Tags:    plan.Tags,
		Systems: plan.Systems,
		Notes:   plan.Notes.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Dns Record",
			"Could not read Id "+fmt.Sprint(dnsRecordId)+": "+err.Error(),
		)
		return
	}

	// update state
	setDnsRecordState(updateDnsRecord, &plan)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState implements tfsdk.Resource
func (dnsRecord) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)

}

func setDnsRecordState(dnsRecord enclaveDns.DnsRecord, state *DnsRecordState) {
	state.Id = types.Int64{Value: int64(dnsRecord.Id)}
	state.Fqdn = types.String{Value: dnsRecord.Fqdn}
}
