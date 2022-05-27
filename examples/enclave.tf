# Super simple example which create a policy along with an enrolment key. It also outputs it so you can see it working!

terraform {
  required_providers {
    enclave = {
      source = "enclave-networks/enclave"
    }
  }
}
variable "enclave_token" {
  type     = string
  nullable = false
}

provider "enclave" {
  token           = var.enclave_token
  organisation_id = "<OrgId>"
}

resource "enclave_policy_acl" "any" {
  protocol = "any"
}

resource "enclave_policy" "testpolicy" {
  description   = "this is a test"
  sender_tags   = ["developer"]
  receiver_tags = ["database"]
  acl = [
    enclave_policy_acl.any,
  ]
}

resource "enclave_enrolment_key" "enrolment" {
  description   = "Enrolment"
  approval_mode = "automatic"
  tags = [
    "example"
  ]
}

# This is a sensative value it's only being output here as an example please be wary with sharing this key
output "enrolment_key" {
  value = enclave_enrolment_key.enrolment.key
  sensitive = true
}
