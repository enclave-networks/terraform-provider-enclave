terraform {
  required_providers {
    enclave = {
      source = "enclave-networks/enclave"
    }
  }
}
variable "enclave_token" {
    type = string
    nullable = false
}

provider "enclave" {
    token = var.enclave_token
    organisation = "Thomas's Org"
}

resource "enclave_policy_acl" "any" {
  protocol = "any"
}

resource "enclave_policy" "testpolicy" {
  description = "this is a test"
  acl = [
    enclave_policy_acl.any,
  ]
}

resource "enclave_enrolment_key" "enrolment" {
  description = "Enrolment"
  approval_mode = "automatic"
  tags = [
    "example"
  ]
}

output "enrolment_key" {
  value = enclave_enrolment_key.enrolment.key
}