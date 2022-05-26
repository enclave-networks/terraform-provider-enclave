package enclave

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnrolmentKeyState struct {
	Id           types.Int64  `tfsdk:"id"`
	Key          types.String `tfsdk:"key"`
	Type         types.String `tfsdk:"type"`
	ApprovalMode types.String `tfsdk:"approval_mode"`
	Description  types.String `tfsdk:"description"`
	Tags         []string     `tfsdk:"tags"`
}

type PolicyState struct {
	Id           types.Int64      `tfsdk:"id"`
	Description  types.String     `tfsdk:"description"`
	Notes        types.String     `tfsdk:"notes"`
	IsEnabled    types.Bool       `tfsdk:"is_enabled"`
	SenderTags   []string         `tfsdk:"sender_tags"`
	RecieverTags []string         `tfsdk:"reciever_tags"`
	Acl          []PolicyAclState `tfsdk:"acl"`
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
}
