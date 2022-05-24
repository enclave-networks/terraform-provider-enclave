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

type dnsZoneResourceType struct{}

func (d dnsZoneResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.Int64Type,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"notes": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

// New resource instance
func (d dnsZoneResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return dnsZone{
		provider: *(p.(*provider)),
	}, nil
}

type dnsZone struct {
	provider provider
}

// Create implements tfsdk.Resource
func (d dnsZone) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !d.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, "+
				"likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan DnsZoneState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dnsZoneCreate := enclaveDns.DnsZoneCreate{
		Name:  plan.Name.Value,
		Notes: plan.Notes.Value,
	}

	// create request
	dnsZone, err := d.provider.client.Dns.CreateZone(dnsZoneCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating DnsZone in enclave",
			err.Error(),
		)
		return
	}

	setDnsZoneStateId(dnsZone, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements tfsdk.Resource
func (d dnsZone) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// read state
	var state DnsZoneState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dnsZoneId := enclaveDns.DnsZoneId(state.Id.Value)

	//call api to delete
	_, err := d.provider.client.Dns.DeleteZone(dnsZoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Dns Zone",
			"Could not read Id "+fmt.Sprint(dnsZoneId)+": "+err.Error(),
		)
		return
	}

	// remove resource
	resp.State.RemoveResource(ctx)
}

// Read implements tfsdk.Resource
func (d dnsZone) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state DnsZoneState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dnsZoneId := enclaveDns.DnsZoneId(state.Id.Value)

	dnsZone, err := d.provider.client.Dns.GetZone(dnsZoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Dns Zone",
			"Could not read Id "+fmt.Sprint(dnsZoneId)+": "+err.Error(),
		)
		return
	}

	setDnsZoneStateId(dnsZone, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update implements tfsdk.Resource
func (d dnsZone) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// read state to get Id
	var state DnsZoneState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get executing plan updates
	var plan DnsZoneState
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	dnsZoneId := enclaveDns.DnsZoneId(state.Id.Value)

	updateDnsZone, err := d.provider.client.Dns.UpdateZone(dnsZoneId, enclaveDns.DnsZonePatch{
		Notes: plan.Notes.Value,
		Name:  plan.Name.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Dns Zone",
			"Could not read Id "+fmt.Sprint(dnsZoneId)+": "+err.Error(),
		)
		return
	}

	// update state
	setDnsZoneStateId(updateDnsZone, &plan)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState implements tfsdk.Resource
func (dnsZone) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)

}

func setDnsZoneStateId(dnsZone enclaveDns.DnsZone, state *DnsZoneState) {
	state.Id = types.Int64{Value: int64(dnsZone.Id)}
}
