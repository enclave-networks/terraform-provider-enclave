package enclave

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnrolmentKeyState struct {
	Id                           types.Int64  `tfsdk:"id"`
	Key                          types.String `tfsdk:"key"`
	Type                         types.String `tfsdk:"type"`
	ApprovalMode                 types.String `tfsdk:"approval_mode"`
	Description                  types.String `tfsdk:"description"`
	DisconnectedRetentionMinutes types.Int64  `tfsdk:"disconnected_retention_minutes"`
	Tags                         []string     `tfsdk:"tags"`
}

type PolicyState struct {
	Id                types.Int64      `tfsdk:"id"`
	Description       types.String     `tfsdk:"description"`
	Notes             types.String     `tfsdk:"notes"`
	IsEnabled         types.Bool       `tfsdk:"is_enabled"`
	SenderTags        []string         `tfsdk:"sender_tags"`
	ReceiverTags      []string         `tfsdk:"receiver_tags"`
	Acl               []PolicyAclState `tfsdk:"acl"`
	TrustRequirements []types.Int64    `tfsdk:"trust_requirements"`
}

type PolicyAclState struct {
	Protocol    types.String `tfsdk:"protocol"`
	Ports       types.String `tfsdk:"ports"`
	Description types.String `tfsdk:"description"`
}

type DnsZoneState struct {
	Id    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Notes types.String `tfsdk:"notes"`
}

type DnsRecordState struct {
	Id      types.Int64  `tfsdk:"id"`
	ZoneId  types.Int64  `tfsdk:"zone_id"`
	Name    types.String `tfsdk:"name"`
	Tags    []string     `tfsdk:"tags"`
	Systems []string     `tfsdk:"systems"`
	Notes   types.String `tfsdk:"notes"`
	Fqdn    types.String `tfsdk:"fqdn"`
}

type TrustRequirementState struct {
	Id                 types.Int64             `tfsdk:"id"`
	Description        types.String            `tfsdk:"description"`
	Notes              types.String            `tfsdk:"notes"`
	UserAuthentication UserAuthenticationState `tfsdk:"user_authentication"`
}

type UserAuthenticationState struct {
	Authority     types.String                        `tfsdk:"authority"`
	AzureTenantId types.String                        `tfsdk:"azure_tenant_id"`
	AzureGroupId  types.String                        `tfsdk:"azure_group_id"`
	Mfa           types.Bool                          `tfsdk:"mfa"`
	CustomClaims  []TrustRequirementCustomClaimsState `tfsdk:"custom_claims"`
}

type TrustRequirementCustomClaimsState struct {
	Claim types.String `tfsdk:"claim"`
	Value types.String `tfsdk:"value"`
}

type TagState struct {
	Ref               types.String  `tfsdk:"ref"`
	Name              types.String  `tfsdk:"name"`
	Colour            types.String  `tfsdk:"colour"`
	Notes             types.String  `tfsdk:"notes"`
	TrustRequirements []types.Int64 `tfsdk:"trust_requirements"`
}
