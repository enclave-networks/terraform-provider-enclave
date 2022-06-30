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
resource "enclave_dns_zone" "zone1" {
  name = "internal"
}

resource "enclave_dns_record" "record1"{
  name = "terraform-test"
  zone_id = enclave_dns_zone.zone1.id
  tags = [
      "tag1",
      "tag2-dns"
  ]
}
```

## Schema

- `zone_id` - (Optional) A DNS Zone ID which can be retrieved from the resource.

- `name` - (Required) The DNS Record name which also forms the `FQDN`.

- `tags` - (Optional) The list of Tags that this Record will apply to.

- `systems` - (Optional) A list of system IDs this Record will apply to.

- `notes` - (Optional) Notes about this DNS Record.
