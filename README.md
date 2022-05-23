# terraform-provider-enclave
Terraform provider for managing Enclave


# Local Testing
run `go build` to generate an executable for terraform to use

then add the following to `%APPDATA%/terraform.rc` replacing `<PATH>` with the path to this repo such as `C:\\git\\enclave\\terraform-provider-enclave`

```
provider_installation {

  dev_overrides {
      "enclave-networks/enclave" = "<PATH>"
      }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```