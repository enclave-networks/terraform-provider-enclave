---
page_title: "enclave_policy resource - Enclave"
subcategory: ""
description: |-
A policy defines which systems can talk to each other, and what traffic can flow between them.
---

# Resource `enclave_policy`

The policy resource is used to create an enclave policy; this defines how systems communicate, with attributes such as sender/receiver tags and ACLs.

## Example

```terraform
# This policy allows all systems tagged with "dev" or "tester" 
# to access all systems tagged with "server", but not the other 
# way round.
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

## Schema

- `description` - (Required) A brief description of the policy e:g `Development Access`.

- `is_enabled` - (Optional) Is the policy enabled? This defaults to `true`.

- `sender_tags` - (Optional) A list of sender tags to apply to this policy. All systems with this tag will be able to send to the receivers.

- `receiver_tags` - (Optional) A list of receiver tags to apply to this policy. All systems with this tag will be able to receive traffic from the senders.

- `acl` - (Optional) More info can be found in the `policy_acl` section of these docs. If no ACLs are specified, no traffic will flow across the policy.

- `notes` - (Optional) Some notes about the policy.
