package main

import (
	"context"

	terraformEnclave "github.com/enclave-networks/terraform-provider-enclave/enclave"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {

	tfsdk.Serve(context.Background(), terraformEnclave.New, tfsdk.ServeOpts{
		Name: "enclave",
	})
}
