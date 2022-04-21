terraform {
  required_providers {
    enclave = {
      source = "enclave-networks/enclave"
    }
  }
}

# Token can this be provided on apply rather than here? Digital oceon does this
provider "enclave" {
    token = "p9rcFksNsHALkfyqyfgRzYq4AXwcuxr22CN9Mc5PG42umHPUiPhnzX7kiRfdWM3"
    url = "http://localhost:8081/"
}

# Can we use id of resource to specify description
resource "enclave_enrolment_key" "developers" {
    type = "general"
    approval_mode = "automatic"
    description = "developers"
    tags = [
        "tag1",
        "tag2"
    ]
}