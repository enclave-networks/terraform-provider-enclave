package enclave

import (
	"context"
	"fmt"
	"strings"

	enclaveData "github.com/enclave-networks/go-enclaveapi/data"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Can probably use a data source for ACLs however need to understand that more https://www.terraform.io/plugin/framework/data-sources
// https://learn.hashicorp.com/tutorials/terraform/plugin-framework-create?in=terraform/providers

type enrolmentKeyResourceType struct{}

func (e enrolmentKeyResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"type": {
				Type:     types.StringType,
				Required: true,
			},
			"approvalMode": {
				Type:     types.StringType,
				Optional: true,
			},
			"description": {
				Type:     types.StringType,
				Required: true,
			},
			"tags": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
		},
	}, nil
}

// New resource instance
func (e enrolmentKeyResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return enrolmentKey{
		provider: *(p.(*provider)),
	}, nil
}

type enrolmentKey struct {
	provider provider
}

// Create a new resource
func (e enrolmentKey) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !e.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, "+
				"likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from enrolmentKeyAttributes
	var enrolmentKeyAttributes EnrolmentKeyState
	diags := req.Plan.Get(ctx, &enrolmentKeyAttributes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	enrolmentKeyType, err := getType(enrolmentKeyAttributes.Type)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting string to enum for enrolmentKeyType",
			err.Error(),
		)
		return
	}

	approvalModeType, err := getApprovalMode(enrolmentKeyAttributes.ApprovalMode)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting string to enum for approvalModeType",
			err.Error(),
		)
		return
	}

	// create request
	enrolmentKeyResponse, err := e.provider.client.EnrolmentKey.Create(enclaveData.EnrolmentKeyCreate{
		Type:         enrolmentKeyType,
		ApprovalMode: approvalModeType,
		Description:  enrolmentKeyAttributes.Description,
		Tags:         enrolmentKeyAttributes.Tags,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating EnrolmentKey in enclave",
			err.Error(),
		)
		return
	}

	enrolmentKeyState := getStateFromResponse(enrolmentKeyResponse)
	diags = resp.State.Set(ctx, enrolmentKeyState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (e enrolmentKey) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state EnrolmentKeyState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	enrolmentKeyId := state.Id

	currentEnrolmentKey, err := e.provider.client.EnrolmentKey.Get(enrolmentKeyId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading enrolment Key",
			"Could not read Id "+enrolmentKeyId+": "+err.Error(),
		)
		return
	}

	enrolmentKeyState := getStateFromResponse(currentEnrolmentKey)
	diags = resp.State.Set(ctx, enrolmentKeyState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (e enrolmentKey) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// read state
	var state EnrolmentKeyState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	enrolmentKeyId := state.Id

	// call api to update
	updateEnrolmentKey, err := e.provider.client.EnrolmentKey.Update(enrolmentKeyId, enclaveData.EnrolmentKeyPatch{
		Description:  state.Description,
		ApprovalMode: state.ApprovalMode,
		Tags:         state.Tags,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating enrolment Key",
			"Could not read Id "+enrolmentKeyId+": "+err.Error(),
		)
		return
	}

	// update state
	enrolmentKeyState := getStateFromResponse(updateEnrolmentKey)
	diags = resp.State.Set(ctx, enrolmentKeyState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (e enrolmentKey) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// read state
	var state EnrolmentKeyState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	enrolmentKeyId := state.Id

	//call api to delete
	_, err := e.provider.client.EnrolmentKey.Disable(enrolmentKeyId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting enrolment Key",
			"Could not read Id "+enrolmentKeyId+": "+err.Error(),
		)
		return
	}

	// remove resource
	resp.State.RemoveResource(ctx)
}

// Import resource
func (e enrolmentKey) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

// Get EnrolmentKeyType from string
func getType(typeString string) (enclaveData.EnrolmentKeyType, error) {
	switch strings.ToLower(typeString) {
	case "general":
		return enclaveData.GeneralPurpose, nil
	case "ephemeral":
		return enclaveData.Ephemeral, nil
	}

	return "", fmt.Errorf("error when converting %s to EnrolmentKeyType", typeString)
}

//Get EnrolmentKeyApprovalMode from string
func getApprovalMode(approvalModeString string) (enclaveData.EnrolmentKeyApprovalMode, error) {
	switch strings.ToLower(approvalModeString) {
	case "automatic":
		return enclaveData.Automatic, nil
	case "manual":
		return enclaveData.Manual, nil
	}
	return "", fmt.Errorf("error when converting %s to EnrolmentKeyApprovalMode", approvalModeString)
}

func getStateFromResponse(enrolmentKey enclaveData.EnrolmentKey) EnrolmentKeyState {
	var tags []string
	for _, x := range enrolmentKey.Tags {
		tags = append(tags, x.Tag)
	}

	return EnrolmentKeyState{
		Id:           enrolmentKey.Id,
		Type:         (string)(enrolmentKey.Type),
		ApprovalMode: (string)(enrolmentKey.ApprovalMode),
		Description:  enrolmentKey.Description,
		Tags:         tags,
	}
}
