package enclave

import (
	"context"

	enclavePolicy "github.com/enclave-networks/go-enclaveapi/data/policy"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type policyResourceType struct{}

func (p policyResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.Int64Type,
				Computed: true,
			},
			"description": {
				Type:     types.StringType,
				Required: true,
			},
			"notes": {
				Type:     types.StringType,
				Required: false,
			},
			"is_enabled": {
				Type:     types.BoolType,
				Required: false,
			},
			"sender_tags": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"reciever_tags": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
		},
	}, nil
}

func (p policyResourceType) NewResource(_ context.Context, pr tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return policy{
		provider: *(pr.(*provider)),
	}, nil
}

type policy struct {
	provider provider
}

func (p policy) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !p.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, "+
				"likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan PolicyState

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var isEnabled bool
	if plan.IsEnabled.Null {
		isEnabled = true
	}

	policyCreate := enclavePolicy.PolicyCreate{
		Description:  plan.Description.Value,
		IsEnabled:    isEnabled,
		Notes:        plan.Notes.Value,
		SenderTags:   plan.SenderTags,
		RecieverTags: plan.RecieverTags,
	}

	// create request
	policyResponse, err := p.provider.client.PolicyClient.Create(policyCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating EnrolmentKey in enclave",
			err.Error(),
		)
		return
	}

	setPolicyStateId(policyResponse, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p policy) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

}

func (p policy) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

func (p policy) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
}

// Import resource
func (p policy) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

func setPolicyStateId(policy enclavePolicy.Policy, state *PolicyState) {
	state.Id = types.Int64{Value: int64(policy.Id)}
}
