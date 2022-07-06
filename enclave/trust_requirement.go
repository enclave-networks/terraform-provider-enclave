package enclave

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	enclaveTrustRequirement "github.com/enclave-networks/go-enclaveapi/data/trustrequirement"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type trustRequirementResourceType struct{}

func (t trustRequirementResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Optional: true,
			},
			"user_authentication": {
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"authority": {
						Type:     types.StringType,
						Required: true,
					},
					"azure_tenant_id": {
						Type:     types.StringType,
						Optional: true,
					},
					"azure_group_id": {
						Type:     types.StringType,
						Optional: true,
					},
					"mfa": {
						Type:     types.BoolType,
						Optional: true,
					},
					"custom_claims": {
						Type: types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"claim": types.StringType,
									"value": types.StringType,
								},
							},
						},
						Optional: true,
					},
				}),
				Optional: true,
			},
		},
	}, nil
}

func (t trustRequirementResourceType) NewResource(_ context.Context, pr tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return trustRequirement{
		provider: *(pr.(*provider)),
	}, nil
}

type trustRequirement struct {
	provider provider
}

// Create implements tfsdk.Resource
func (t trustRequirement) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !t.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, "+
				"likely because it depends on an unknown value from another resource. "+
				"This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan TrustRequirementState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Let's check lengths and add some warnings
	validateTrustRequirement(plan, &resp.Diagnostics)

	trustRequirementType, config, conditions, err := getTrustRequirementSettings(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating api request",
			err.Error(),
		)
	}

	trustRequirementCreate := enclaveTrustRequirement.TrustRequirementCreate{
		Description: plan.Description.Value,
		Notes:       plan.Notes.Value,
		Type:        trustRequirementType,
		Settings: enclaveTrustRequirement.TrustRequirementSettings{
			Configuration: config,
			Conditions:    conditions,
		},
	}

	// create request
	trustRequirementResponse, err := t.provider.client.TrustRequirements.Create(trustRequirementCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Trust Requirement in enclave",
			err.Error(),
		)
		return
	}

	setTrustRequirementStateId(trustRequirementResponse, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements tfsdk.Resource
func (t trustRequirement) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// read state
	var state TrustRequirementState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustRequirementId := enclaveTrustRequirement.TrustRequirementId(state.Id.Value)

	//call api to delete
	_, err := t.provider.client.TrustRequirements.Delete(trustRequirementId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Trust Requirement",
			"Could not read Id "+fmt.Sprint(trustRequirementId)+": "+err.Error(),
		)
		return
	}

	// remove resource
	resp.State.RemoveResource(ctx)
}

// Read implements tfsdk.Resource
func (t trustRequirement) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state TrustRequirementState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustRequirementId := enclaveTrustRequirement.TrustRequirementId(state.Id.Value)

	currentTrustRequirement, err := t.provider.client.TrustRequirements.Get(trustRequirementId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Trust requirement Id",
			"Could not read Id "+fmt.Sprint(trustRequirementId)+": "+err.Error(),
		)
		return
	}

	setTrustRequirementStateId(currentTrustRequirement, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update implements tfsdk.Resource
func (t trustRequirement) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// read state to get Id
	var state TrustRequirementState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get executing plan updates
	var plan TrustRequirementState
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Let's check lengths and add some warnings
	validateTrustRequirement(plan, &resp.Diagnostics)

	trustRequirementId := enclaveTrustRequirement.TrustRequirementId(state.Id.Value)

	_, config, conditions, err := getTrustRequirementSettings(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating api request",
			err.Error(),
		)
	}

	updateTrustRequirement, err := t.provider.client.TrustRequirements.Update(trustRequirementId, enclaveTrustRequirement.TrustRequirementPatch{
		Description: plan.Description.Value,
		Notes:       plan.Notes.Value,
		Settings: enclaveTrustRequirement.TrustRequirementSettings{
			Configuration: config,
			Conditions:    conditions,
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Trust Requirement",
			"Could not read Id "+fmt.Sprint(trustRequirementId)+": "+err.Error(),
		)
		return
	}

	// update state
	setTrustRequirementStateId(updateTrustRequirement, &plan)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func setTrustRequirementStateId(trustRequirement enclaveTrustRequirement.TrustRequirement, state *TrustRequirementState) {
	state.Id = types.Int64{Value: int64(trustRequirement.Id)}
}

func validateTrustRequirement(plan TrustRequirementState, diagnostics *diag.Diagnostics) {
	if strings.ToLower(plan.UserAuthentication.Authority.Value) == string(Portal) && (!plan.UserAuthentication.AzureGroupId.Null || !plan.UserAuthentication.AzureTenantId.Null) {
		diagnostics.AddWarning(
			"Authetication Authority of Portal does not need any additional properties",
			"The portal Authority type only needs the Authority value")
	}
}

func getTrustRequirementSettings(plan TrustRequirementState) (trustRequirementType enclaveTrustRequirement.TrustRequirementType, config map[string]string, conditions []map[string]string, err error) {
	// UserAuthentication has been set use that to create our maps
	if !reflect.DeepEqual(plan.UserAuthentication, UserAuthenticationState{}) {
		authorityLower := strings.ToLower(plan.UserAuthentication.Authority.Value)

		if authorityLower == string(Portal) {
			return enclaveTrustRequirement.UserAuthentication,
				map[string]string{
					"authority": authorityLower,
				},
				[]map[string]string{},
				nil
		}

		if authorityLower == string(Azure) {
			conditions := []map[string]string{
				{
					"claim": "groups",
					"value": plan.UserAuthentication.AzureGroupId.Value,
				},
			}

			// only add mfa claim if set and true
			if !plan.UserAuthentication.Mfa.Null && plan.UserAuthentication.Mfa.Value {
				conditions = append(conditions, map[string]string{
					"claim": "amr",
					"value": "mfa",
				})
			}

			if len(plan.UserAuthentication.CustomClaims) > 0 {
				customClaims := make([]map[string]string, len(plan.UserAuthentication.CustomClaims))
				for i, item := range plan.UserAuthentication.CustomClaims {
					customClaims[i] = map[string]string{
						"claim": item.Claim.Value,
						"value": item.Value.Value,
					}
				}

				conditions = append(conditions, customClaims...)
			}

			return enclaveTrustRequirement.UserAuthentication,
				map[string]string{
					"authority": authorityLower,
					"tenantId":  plan.UserAuthentication.AzureTenantId.Value,
				}, conditions, nil
		}

	}

	// We shouldn't ever really get here but just in case we'll inform the user they've not got a value
	return -1,
		map[string]string{},
		[]map[string]string{},
		fmt.Errorf("could not get trust requirement settings please ensure you have a type object created refer to the docs for more information")
}

type TrustRequirementAuthorityType string

const (
	Portal TrustRequirementAuthorityType = "portal"
	Azure  TrustRequirementAuthorityType = "azure"
)
