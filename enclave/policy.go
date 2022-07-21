package enclave

import (
	"context"
	"fmt"
	"strings"

	enclavePolicy "github.com/enclave-networks/go-enclaveapi/data/policy"
	"github.com/hashicorp/terraform-plugin-framework/attr"

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
				Optional: true,
			},
			"is_enabled": {
				Type:     types.BoolType,
				Optional: true,
			},
			"sender_tags": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"receiver_tags": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"trust_requirements": {
				Type: types.ListType{
					ElemType: types.Int64Type,
				},
				Optional: true,
			},
			"acl": {
				Type: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"protocol":    types.StringType,
							"ports":       types.StringType,
							"description": types.StringType,
						},
					},
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

	// Let's check lengths and add some warnings
	checkForPolicyWarnings(plan, &resp.Diagnostics)

	policyAcl, err := toPolicyAcl(plan.Acl)
	if err != nil {
		resp.Diagnostics.AddError("Error converting ACL", err.Error())
	}

	policyCreate := enclavePolicy.PolicyCreate{
		Description:       plan.Description.Value,
		IsEnabled:         isEnabled,
		Notes:             plan.Notes.Value,
		SenderTags:        plan.SenderTags,
		ReceiverTags:      plan.ReceiverTags,
		Acls:              policyAcl,
		TrustRequirements: toTrustRequirementIdArray(plan.TrustRequirements),
	}

	// create request
	policyResponse, err := p.provider.client.Policies.Create(policyCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Policy in enclave",
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
	var state PolicyState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := enclavePolicy.PolicyId(state.Id.Value)

	currentPolicy, err := p.provider.client.Policies.Get(policyId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading policy Id",
			"Could not read Id "+fmt.Sprint(policyId)+": "+err.Error(),
		)
		return
	}

	setPolicyStateId(currentPolicy, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p policy) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// read state to get Id
	var state PolicyState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get executing plan updates
	var plan PolicyState
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Let's check lengths and add some warnings
	checkForPolicyWarnings(plan, &resp.Diagnostics)

	policyAcl, err := toPolicyAcl(plan.Acl)
	if err != nil {
		resp.Diagnostics.AddError("Error converting ACL", err.Error())
	}

	policyId := enclavePolicy.PolicyId(state.Id.Value)

	updatePolicy, err := p.provider.client.Policies.Update(policyId, enclavePolicy.PolicyPatch{
		Description:       plan.Description.Value,
		IsEnabled:         plan.IsEnabled.Value,
		SenderTags:        plan.SenderTags,
		ReceiverTags:      plan.ReceiverTags,
		Notes:             plan.Notes.Value,
		Acls:              policyAcl,
		TrustRequirements: toTrustRequirementIdArray(plan.TrustRequirements),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Policy",
			"Could not read Id "+fmt.Sprint(policyId)+": "+err.Error(),
		)
		return
	}

	// update state
	setPolicyStateId(updatePolicy, &plan)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p policy) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// read state
	var state PolicyState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := enclavePolicy.PolicyId(state.Id.Value)

	//call api to delete
	_, err := p.provider.client.Policies.Delete(policyId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Policy",
			"Could not read Id "+fmt.Sprint(policyId)+": "+err.Error(),
		)
		return
	}

	// remove resource
	resp.State.RemoveResource(ctx)
}

// Import resource
func (p policy) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

func setPolicyStateId(policy enclavePolicy.Policy, state *PolicyState) {
	state.Id = types.Int64{Value: int64(policy.Id)}
}

func checkForPolicyWarnings(plan PolicyState, diagnostics *diag.Diagnostics) {
	if len(plan.Acl) == 0 {
		diagnostics.AddWarning(
			"No ACLs defined",
			"No network traffic is permitted by this policy until at least one acl is applied. "+
				"Allow network traffic for this policy by adding an acl entry for the policy")
	}

	if len(plan.SenderTags) == 0 {
		diagnostics.AddWarning(
			"No Sender Tags defined",
			"Please define a sender tag in order for this Policy to create connectivity")
	}

	if len(plan.ReceiverTags) == 0 {
		diagnostics.AddWarning(
			"No Receiver Tags defined",
			"Please define a receiver tag in order for this Policy to create connectivity")
	}
}

func toPolicyAcl(pacl []PolicyAclState) ([]enclavePolicy.PolicyAcl, error) {

	result := make([]enclavePolicy.PolicyAcl, len(pacl))

	for i, acl := range pacl {
		protcol, err := isValidProtocol(acl.Protocol.Value)

		if err != nil {
			return nil, err
		}

		result[i] = enclavePolicy.PolicyAcl{
			Ports:       acl.Ports.Value,
			Protocol:    protcol,
			Description: acl.Description.Value,
		}
	}

	return result, nil
}

func isValidProtocol(protocol string) (enclavePolicy.PolicyAclProtocol, error) {
	switch strings.ToLower(protocol) {
	case "any":
		return enclavePolicy.Any, nil
	case "tcp":
		return enclavePolicy.Tcp, nil
	case "udp":
		return enclavePolicy.Udp, nil
	case "icmp":
		return enclavePolicy.Icmp, nil
	}

	return "", fmt.Errorf("invalid protcol specified must be one of: any, tcp, udp, icmp")
}
