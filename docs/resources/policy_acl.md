---
page_title: "policy_acl Resource - Enclave"
subcategory: ""
description: |-
A Policy ACL defines what sort of access a system has in relation to the current policy.
---

# Resource `enclave_policy_acl`

This will allow you to configure Access Control Lists, specifying what traffic you want to allow through your policy from a sender to a receiver. This layout also allows for reusing these ACLs.

The enclave Policy ACL needs to be used in conjunction with a Policy as it will not create anything unless it's associated to an `enclave_policy`.

## Example

```terraform
resource "enclave_policy_acl" "any" {
  protocol = "any"
}

resource "enclave_policy_acl" "http" {
  protocol = "tcp"
  ports = "80"
}

resource "enclave_policy" "testpolicy" {
  description = "this is a policy that allows any traffic"
  acl = [
    enclave_policy_acl.any
  ]
}

resource "enclave_policy" "testpolicy2" {
  description = "this policy only allows http traffic"
  acl = [
    enclave_policy_acl.http
  ]
}
```

## Schema

- `protocol` - (Required) The Protocol type can be one of the following `any`, `tcp`, `udp`, `icmp`.

- `ports`- (Optional) A port range or a single port e:g `8000-8080` or `8080`. Only valid for `tcp` and `udp` protocol values. If not set then all traffic of the specified `protocol` is allowed.

- `description` - (Optional) A description of this ACL.
