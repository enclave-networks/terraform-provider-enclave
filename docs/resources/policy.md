---
page_title: "policy Resource - Enclave"
subcategory: ""
description: |-
A Policy is what defines the rules on which systems can talk to each other.
---

# Resource `enclave_policy`

The Policy resource is used to create a policy enclave this defines how systems communicate with attributes such as sender/receiver tags and ACLs.

# Example

```terraform
resource "enclave_policy" "testpolicy" {
  description = "this is a test"
  notes = "i'll use this to show how it works"
  is_enabled = true
  sender_tags = [
    "dev",
    "tester"
  ]
  receiver_tags = [
    "server"
  ]
}
```

# Schema
  - `description` - (Required) A brief description of the policy e:g `Development Access`.
  - `notes` - (Optional) Some notes about the policy.
  - `is_enabled` - (Optional) Is the policy enabled this defaults to `true`.
  - `sender_tags` - (Optional) A list of sender tags to apply to this policy.
  - `receiver_tags` - (Optional) A list of receiver tags to apply to this policy.
  - `acl` - (Optional) More info can be found in the `policy_acl` section of these docs.