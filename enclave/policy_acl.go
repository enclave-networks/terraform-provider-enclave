package enclave

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type policyAclResourceType struct{}

func (pa policyAclResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"protocol": {
				Type:     types.StringType,
				Required: true,
			},
			"ports": {
				Type:     types.StringType,
				Optional: true,
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

type policyAcl struct {
	provider provider
}

// Create implements tfsdk.Resource
func (pa policyAcl) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !pa.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, "+
				"likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan PolicyAclState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements tfsdk.Resource
func (pa policyAcl) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// read state
	var state PolicyAclState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// remove resource
	resp.State.RemoveResource(ctx)
}

// Read implements tfsdk.Resource
func (pa policyAcl) Read(context.Context, tfsdk.ReadResourceRequest, *tfsdk.ReadResourceResponse) {
	//Do Nothing
}

// Update implements tfsdk.Resource
func (pa policyAcl) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !pa.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, "+
				"likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan PolicyAclState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState implements tfsdk.Resource
func (pa policyAcl) ImportState(context.Context, tfsdk.ImportResourceStateRequest, *tfsdk.ImportResourceStateResponse) {
	// Do Nothing
}

func (pa policyAclResourceType) NewResource(_ context.Context, pr tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return policyAcl{
		provider: *(pr.(*provider)),
	}, nil
}
