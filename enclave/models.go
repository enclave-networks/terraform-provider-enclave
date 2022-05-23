package enclave

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnrolmentKeyState struct {
	Id           types.Int64  `tfsdk:"id"`
	Type         types.String `tfsdk:"type"`
	ApprovalMode types.String `tfsdk:"approval_mode"`
	Description  types.String `tfsdk:"description"`
	Tags         []string     `tfsdk:"tags"`
}

type PolicyState struct {
	Id           types.Int64  `tfsdk:"id"`
	Description  types.String `tfsdk:"description"`
	Notes        types.String `tfsdk:"Notes"`
	IsEnabled    types.Bool   `tfsdk:"is_enabled"`
	SenderTags   []string     `tfsdk:"sender_tags"`
	RecieverTags []string     `tfsdk:"reciever_tags"`
}
