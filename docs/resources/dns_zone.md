---
page_title: "dns_zone Resource - Enclave"
subcategory: ""
description: |-
A DNS Zone is a suffix for your DNS Records; by default `.enclave` is available for all plans.
---

# Resource `enclave_dns_zone`

A DNS Zone is a container/suffix for a DNS Record this allows you to create multiple Records within it. For example you can create a Zone called `internal` that will have the `.internal` suffix on all records.

Custom DNS zones are only 

More information can be found in the [enclave docs](https://docs.enclave.io/management/dns/).

# Example

```terraform
resource "enclave_dns_zone" "zone1" {
  name = "internal"
}
```

# Schema

- `name` - (Required) The name without spaces of the DNS Zone. This will them become the Suffix e:g `zone-name` becomes `.zone-name`.

- `notes` - (Optional) Some notes on what this DNS Zone is used for.
