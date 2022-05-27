# terraform-provider-enclave

Welcome to the [Enclave Terraform Provider](https://registry.terraform.io/providers/enclave-networks/enclave/latest) Repo. This provider gives access to some of the more common functions of the [Enclave Api](https://api.enclave.io/index.html). The documentation on the [terraform registry](https://registry.terraform.io/providers/enclave-networks/enclave/latest/docs) will be the most up to date but for a basic getting started please check below!


# Getting Started
First add the following to your `.tf` file and then run `terraform init` to pull down the required files. From here you can provision enclave resources using this [Documentation](https://registry.terraform.io/providers/enclave-networks/enclave/latest/docs) as a guide!
```terraform
terraform {
  required_providers {
    enclave = {
      source = "enclave-networks/enclave"
    }
  }
}

variable "enclave_token" {
    type = string
    nullable = false
}


provider "enclave" {
    token = var.enclave_token
}
```





# Contributing
## Issues/Suggestions
If you have any issues please provide log messages if possible and head over to the Issues tab to create an issue for us. Same can be said for features but don't worry you won't need any logs!

## Local Testing
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