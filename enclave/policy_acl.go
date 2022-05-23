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
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Required: false,
			},
		},
	}, nil
}

type policyAcl struct {
	provider provider
}

// Create implements tfsdk.Resource
func (policyAcl) Create(context.Context, tfsdk.CreateResourceRequest, *tfsdk.CreateResourceResponse) {
	panic("unimplemented")
}

// Delete implements tfsdk.Resource
func (policyAcl) Delete(context.Context, tfsdk.DeleteResourceRequest, *tfsdk.DeleteResourceResponse) {
	panic("unimplemented")
}

// ImportState implements tfsdk.Resource
func (policyAcl) ImportState(context.Context, tfsdk.ImportResourceStateRequest, *tfsdk.ImportResourceStateResponse) {
	panic("unimplemented")
}

// Read implements tfsdk.Resource
func (policyAcl) Read(context.Context, tfsdk.ReadResourceRequest, *tfsdk.ReadResourceResponse) {
	panic("unimplemented")
}

// Update implements tfsdk.Resource
func (policyAcl) Update(context.Context, tfsdk.UpdateResourceRequest, *tfsdk.UpdateResourceResponse) {
	panic("unimplemented")
}

func (pa policyAclResourceType) NewResource(_ context.Context, pr tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return policyAcl{
		provider: *(pr.(*provider)),
	}, nil
}
