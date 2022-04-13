package enclave

import (
	"context"

	enclaveData "github.com/enclave-networks/go-enclaveapi/data"
	"github.com/enclave-networks/go-enclaveapi/enclave"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *enclave.OrganisationClient
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"token": {
				Type:     types.StringType,
				Required: true,
			},
			"organisation": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

type providerData struct {
	Token        types.String `tfsdk:"token"`
	Organisation types.String `tfsdk:"organisation"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var token string
	if config.Token.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as token",
		)
		return
	}

	if !config.Token.Null {
		token = config.Token.Value
	}

	if token == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find token",
			"Token cannot be an empty string",
		)
		return
	}

	var organisation string
	if config.Organisation.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as organisation",
		)
		return
	}

	if !config.Token.Null {
		organisation = config.Organisation.Value
	}

	c := enclave.CreateClient(&token)
	orgs, err := c.GetOrgs()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting enclave orgs",
			"Please ensure you have a valid Token",
		)
		return
	}

	var currentOrg enclaveData.AccountOrganisation
	if &organisation == nil {
		currentOrg = *orgs[0]
	} else {
		for _, o := range orgs {
			if o.OrgName == &organisation {
				currentOrg = *o
			}
		}
	}

	if &currentOrg == nil {
		resp.Diagnostics.AddError(
			"Could not find org",
			"Please ensure you have specified the correct org and you have access to it",
		)
		return
	}

	p.client = c.CreateOrganisationClient(&currentOrg)
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"enclave_enrolment_key": enrolmentKeyResourceType{},
		// Add more resource types here
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
