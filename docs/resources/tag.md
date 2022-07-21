---
page_title: "enclave_tag resource - Enclave"
subcategory: ""
description: |-
Tags are free-form text labels attached to one or more systems in your account which allow administrators to group systems together which share similar characteristics (i.e. business unit, security level or function) in order to build connectivity between them.
---

# Resource `enclave_tag`

A tag can either be created as a text field on other resources such as policies or enrolment keys or using this resource. Here you can assign Trust Requirements and Colours to your tags.

# Example

```terraform
resource "enclave_tag" "tag_1" {
  name = "this-is-a-tag"
  colour = "#fcba03"
  notes = "Here are some notes"
  trust_requirements = [
    enclave_trust_requirement.my_first_trust.id
  ]
}
```

## Schema

- `name` (Required) The name of the Tag.

- `colour` (Optional) A Hex Colour Code for the Tag that will be used to display the Tag in the [portal](https://portal.enclave.io/).

- `notes` - (Optional) Some notes on what this Tag is used for.

- `trust_requirements` (Optional) An array of Trust Requirement IDs that will apply to this Tag before connectivity is established.
