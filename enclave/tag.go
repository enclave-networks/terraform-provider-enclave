package enclave

import (
	"context"

	enclaveTag "github.com/enclave-networks/go-enclaveapi/data/tag"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type tagResourceType struct{}

func (t tagResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"ref": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"colour": {
				Type:     types.StringType,
				Optional: true,
			},
			"notes": {
				Type:     types.StringType,
				Optional: true,
			},
			"trust_requirements": {
				Type: types.ListType{
					ElemType: types.Int64Type,
				},
				Optional: true,
			},
		},
	}, nil
}

func (t tagResourceType) NewResource(_ context.Context, pr tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return tag{
		provider: *(pr.(*provider)),
	}, nil
}

type tag struct {
	provider provider
}

func (t tag) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !t.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, "+
				"likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan TagState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagCreate := enclaveTag.TagCreate{
		Tag:               plan.Name.Value,
		Colour:            plan.Colour.Value,
		Notes:             plan.Notes.Value,
		TrustRequirements: toTrustRequirementSlice(plan.TrustRequirements),
	}

	// create request
	tagResponse, err := t.provider.client.Tags.Create(tagCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Tag in enclave",
			err.Error(),
		)
		return
	}

	setTagRef(tagResponse, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (t tag) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state TagState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	currentTag, err := t.provider.client.Tags.Get(state.Ref.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tag ref",
			"Could not read Id "+state.Ref.Value+": "+err.Error(),
		)
		return
	}

	setTagRef(currentTag, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (t tag) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// read state to get Id
	var state TagState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get executing plan updates
	var plan TagState
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	updatedTag, err := t.provider.client.Tags.Update(state.Ref.Value, enclaveTag.TagPatch{
		Tag:               plan.Name.Value,
		Colour:            plan.Colour.Value,
		Notes:             plan.Notes.Value,
		TrustRequirements: toTrustRequirementSlice(plan.TrustRequirements),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Tag",
			"Could not read Id "+state.Ref.Value+": "+err.Error(),
		)
		return
	}

	// update state
	setTagRef(updatedTag, &plan)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (t tag) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// read state
	var state TagState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//call api to delete
	_, err := t.provider.client.Tags.Delete(state.Ref.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Tag",
			"Could not read Id "+state.Ref.Value+": "+err.Error(),
		)
		return
	}

	// remove resource
	resp.State.RemoveResource(ctx)
}

// Import resource
func (t tag) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

func setTagRef(tag enclaveTag.DetailedTag, state *TagState) {
	state.Ref = types.String{Value: string(tag.Ref)}
}
