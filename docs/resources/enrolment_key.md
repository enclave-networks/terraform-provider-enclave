---
page_title: "enrolment_key Resource - Enclave"
subcategory: ""
description: |-
An enrolment key allows you to enrol a system into your Enclave organisation.
---

# Resource `enclave_enrolment_key`

This Enrolment Key resource allows you to generate an Enrolment Key as well as allowing you to access and output that key.

## Example

```terraform
resource "enclave_enrolment_key" "keyname" {
    type = "general"
    approval_mode = "automatic"
    description = "this is a description"
    tags = [
        "tag1",
        "tag2"
    ]
}

output "key_value" {
    value = enclave_enrolment_key.keyname.key
}
```

## Schema

- `type` - (Optional) Can be either `general` or `ephemeral`. Systems enrolled with an `ephemeral` key are removed from your Enclave organisation when they stop, whereas systems enrolled with a `general` key remain in the account. Consider using `ephemeral` keys for containers, kubernetes, etc.

- `approval_mode` - (Optional) Can be either `automatic` or `manual` Will default to manual if not set.

- `description` - (Required) A description of the Enrolment Key.

- `tags` - (Optional) An array of tags that will automatically be applied to any system enrolled with this key.

## Attributes

The following additional attributes are available for all keys:

- `key` - This is the Enrolment Key that is generated after a successful API request.
