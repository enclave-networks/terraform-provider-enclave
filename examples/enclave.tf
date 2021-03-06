# Enclave Terraform Example - Developer Database access
#
# Super simple example which creates a policy to allow developers 
# access to SQL databases, along with an enrolment key for the databases.
# It also outputs the key so you can see it working!

terraform {
  required_providers {
    enclave = {
      source = "enclave-networks/enclave"
    }
  }
}

variable "enclave_token" {
  type      = string
  nullable  = false
  sensitive = true
}

provider "enclave" {
  token           = var.enclave_token
  organisation_id = "<orgId>"
}

resource "enclave_policy_acl" "sql" {
  protocol = "tcp"
  ports = "1433"
}

resource "enclave_policy" "devs_to_db" {
  description   = "Allow devs to access the database"
  sender_tags   = ["developer"]
  receiver_tags = ["database"]
  acl = [
    enclave_policy_acl.sql
  ]
}

resource "enclave_enrolment_key" "db_enrolment" {
  description   = "Enrolment"
  approval_mode = "automatic"
  tags = [
    "database"
  ]
}

resource "enclave_tag" "tag_1" {
  name = "this-is-a-tag"
  trust_requirements = [
    enclave_trust_requirement.my_first_trust.id
  ]
}

resource "enclave_trust_requirement" "my_first_trust" {
  description = "Azure Access"
  user_authentication = {
    authority = "Azure" 
    azure_tenant_id = "<tenant-id>"
    azure_group_id = "<group-id>"
    mfa = true
    custom_claims = [
      {
        claim = "<claim-name>"
        value = "<value>"
      }
    ]
  }
}

# This is a sensitive value; it's only being output here as an example please be wary with sharing this key
output "enrolment_key" {
  value = enclave_enrolment_key.db_enrolment.key
  sensitive = true
}
