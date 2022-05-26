package enclave

import (
	"context"
	"fmt"
	"strings"

	enclaveEnrolmentKey "github.com/enclave-networks/go-enclaveapi/data/enrolmentkey"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Can probably use a data source for ACLs however need to understand that more https://www.terraform.io/plugin/framework/data-sources
// https://learn.hashicorp.com/tutorials/terraform/plugin-framework-create?in=terraform/providers
//https://www.terraform.io/plugin/framework/resources

type enrolmentKeyResourceType struct{}

func (e enrolmentKeyResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.Int64Type,
				Computed: true,
			},
			"key": {
				Type:     types.StringType,
				Computed: true,
			},
			"type": {
				Type:     types.StringType,
				Optional: true,
			},
			"approval_mode": {
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

	// Retrieve values from plan
	var plan EnrolmentKeyState

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var enrolmentKeyType enclaveEnrolmentKey.EnrolmentKeyType
	if !plan.Type.Null {
		val, err := getType(plan.Type.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting string to enum for enrolmentKeyType",
				err.Error(),
			)
			return
		}

		enrolmentKeyType = val
	} else {
		val := enclaveEnrolmentKey.GeneralPurpose
		enrolmentKeyType = val
	}

	var approvalModeType enclaveEnrolmentKey.EnrolmentKeyApprovalMode
	if !plan.ApprovalMode.Null {
		val, err := getApprovalMode(plan.ApprovalMode.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting string to enum for approvalModeType",
				err.Error(),
			)
			return
		}

		approvalModeType = val
	} else {
		val := enclaveEnrolmentKey.Manual
		approvalModeType = val
	}

	enrolmentKeyCreate := enclaveEnrolmentKey.EnrolmentKeyCreate{
		Type:         enrolmentKeyType,
		ApprovalMode: approvalModeType,
		Description:  plan.Description.Value,
		Tags:         plan.Tags,
	}

	// create request
	enrolmentKeyResponse, err := e.provider.client.EnrolmentKey.Create(enrolmentKeyCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating EnrolmentKey in enclave",
			err.Error(),
		)
		return
	}

	setEnrolmentKeyStateValues(enrolmentKeyResponse, &plan)

	diags = resp.State.Set(ctx, plan)
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

	enrolmentKeyId := enclaveEnrolmentKey.EnrolmentKeyId(state.Id.Value)

	currentEnrolmentKey, err := e.provider.client.EnrolmentKey.Get(enrolmentKeyId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading enrolment Key",
			"Could not read Id "+fmt.Sprint(enrolmentKeyId)+": "+err.Error(),
		)
		return
	}

	setEnrolmentKeyStateValues(currentEnrolmentKey, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (e enrolmentKey) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// read state to get Id
	var state EnrolmentKeyState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get executing plan updates
	var plan EnrolmentKeyState
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	enrolmentKeyId := enclaveEnrolmentKey.EnrolmentKeyId(state.Id.Value)

	var approvalModeType enclaveEnrolmentKey.EnrolmentKeyApprovalMode
	if !plan.ApprovalMode.Null {
		val, err := getApprovalMode(plan.ApprovalMode.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting string to enum for approvalModeType",
				err.Error(),
			)
			return
		}

		approvalModeType = val
	} else {
		val := enclaveEnrolmentKey.Manual
		approvalModeType = val
	}

	// call api to update
	updateEnrolmentKey, err := e.provider.client.EnrolmentKey.Update(enrolmentKeyId, enclaveEnrolmentKey.EnrolmentKeyPatch{
		Description:  plan.Description.Value,
		ApprovalMode: approvalModeType,
		Tags:         plan.Tags,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating enrolment Key",
			"Could not read Id "+fmt.Sprint(enrolmentKeyId)+": "+err.Error(),
		)
		return
	}

	// update state
	setEnrolmentKeyStateValues(updateEnrolmentKey, &plan)

	diags = resp.State.Set(ctx, plan)
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
	_, err := e.provider.client.EnrolmentKey.Disable(int(enrolmentKeyId.Value))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting enrolment Key",
			"Could not read Id "+fmt.Sprint(enrolmentKeyId.Value)+": "+err.Error(),
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
func getType(typeString string) (enclaveEnrolmentKey.EnrolmentKeyType, error) {
	switch strings.ToLower(typeString) {
	case "general":
		return enclaveEnrolmentKey.GeneralPurpose, nil
	case "ephemeral":
		return enclaveEnrolmentKey.Ephemeral, nil
	}

	return "", fmt.Errorf("error when converting %s to EnrolmentKeyType", typeString)
}

//Get EnrolmentKeyApprovalMode from string
func getApprovalMode(approvalModeString string) (enclaveEnrolmentKey.EnrolmentKeyApprovalMode, error) {
	switch strings.ToLower(approvalModeString) {
	case "automatic":
		return enclaveEnrolmentKey.Automatic, nil
	case "manual":
		return enclaveEnrolmentKey.Manual, nil
	}
	return "", fmt.Errorf("error when converting %s to EnrolmentKeyApprovalMode", approvalModeString)
}

func setEnrolmentKeyStateValues(enrolmentKey enclaveEnrolmentKey.EnrolmentKey, state *EnrolmentKeyState) {
	state.Id = types.Int64{Value: int64(enrolmentKey.Id)}
	state.Key = types.String{Value: enrolmentKey.Key}
}
