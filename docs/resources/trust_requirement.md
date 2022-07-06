---
page_title: "enclave_trust_requirement resource - Enclave"
subcategory: ""
description: |-
A Trust Requirement can be applied to tags and policies to allow you to apply additional trust checks to systems which are enforced prior to any connectivity being established.
---

# Resource `enclave_trust_requirement`

A Trust Requirement can be created and attached to either a Policy or Tag this defines a set of additional trust checks that are enforced prior to any connectivity being established

## Example

```terraform
resource "enclave_trust_requirement" "my_first_trust" {
  description = "Azure Access"
  user_authentication = {
    authority = "Azure" 
    azure_tenant_id = "<tenant-id>"
    azure_group_id = "<group-id>"
    mfa = true
    custom_claims = [
      {
        claim = "<claim-name>"
        value = "<value>"
      }
    ]
  }
}
```

## Schema

- `description` - (Required) A description of the Trust Requirement.

- `user_authentication` - (Optional) An object used to define a User Authentication Trust Requirement.

  - `authority` - (Required) The type of authority currently only `Portal` and `Azure` are supported.

  - `azure_tenant_id` - (Required if Azure Authority Specified) The azure tenant ID.

  - `azure_group_id` - (Optional) An Azure Group ID.

  - `mfa` - (Optional) Require Multi Factor Authentication.

  - `custom_claims` (Optional) A list of custom claims.
    
    - `claim` (Required) The Name of the custom claim.

    - `value` (Required) The Value of the custom claim.
