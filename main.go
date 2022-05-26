package main

import (
	"context"

	terraformEnclave "github.com/enclave-networks/terraform-provider-enclave/enclave"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), terraformEnclave.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/enclave-networks/enclave",
	})
}
