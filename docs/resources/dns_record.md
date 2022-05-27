---
page_title: "dns_record Resource - Enclave"
subcategory: ""
description: |-
A DNS Record defines a DNS Name to systems or tags
---

# Resource `enclave_dns_record`

A DNS Record can be created and attached to both a system or a tag; it's recommended to attach it to a tag as it's far less brittle. More information can be found on the [enclave docs](https://docs.enclave.io/management/dns/#adding-a-dns-record)

A Record can also to be attached to a custom `dns_zone`. If no `zone_id` is specified it'll be attached to the default `enclave` zone

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
