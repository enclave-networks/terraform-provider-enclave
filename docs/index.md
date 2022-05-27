---
page_title: "Provider: Enclave"
subcategory: ""
description: |-
A Terraform Provider for working with the Enclave API.
---

# Enclave Provider

This provider gives access to the main configuration aspects of Enclave. Use the navigation to the left to read about the available resources.

## Examples

Create your API Token [here](https://portal.enclave.io/account).

Terraform variables should be used when storing your API Token

```terraform
provider "enclave" {
    token = "my-token"
}
```

Variable Example
```terraform
variable "enclave_token" {
    type = string
    nullable = false
    sensitive = true
}


provider "enclave" {
    token = var.enclave_token
}
```

This variable can then be set through the [Terraform CLI](https://www.terraform.io/cli):

```bash
terraform apply -var="enclave_token=my-token"
```

You can also do it through an environment variable:

```bash
export TF_VAR_enclave_token=my-token
terraform apply
```

More examples can be found in our [github repo](https://github.com/enclave-networks/terraform-provider-enclave).

## Schema

### Required

- `token` - string - Your API token from [here](https://portal.enclave.io/account).

### Optional

- `organisation_id` - (Optional unless you have more than one Org associated to your `token`) The ID of your organisation which can be found in the [settings page](https://portal.enclave.io/my/settings), but can also be found in the portal URL (e.g. `https://portal.enclave.io/org/<orgId>/systems`). Organisation IDs are not secret, and can freely be shared. Only users with access to that organisation will be able to invoke APIs against it.

- `url` - string - The Base API Url leave this blank to use the default of `https://api.enclave.io`
