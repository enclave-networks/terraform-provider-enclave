---
page_title: "enrolment_key Resource - Enclave"
subcategory: ""
description: |-
An enrolment key allows you to enrol a system.
---

# Resource `enclave_enrolment_key`
This Enrolment Key resource allows you to generate an Enrolment Key as well as allowing you to access and output that key.


# Example
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

# Schema
- `type` - Can be either `ephemeral` or `general` this will set the state of the enrolled system. Will default to general if not set
- `approval_mode` - Can be either `automatic` or `manual` Will default to manual if not set
- `description` - A description of the Enrolment Key 
- `tags` - An array of tags that will automatically be applied to any system enrolled with this key

# Outputs
- `key` - This is the Enrolment Key that is generated after a successful API request