terraform {
  required_providers {
    enclave = {
      source = "enclave-networks/enclave"
    }
  }
}

# Token can this be provided on apply rather than here? Digital oceon does this
provider "enclave" {
    token = "my-token"
}

resource "enclave_policy_acl" "any" {
  protocol = "any"
}

resource "enclave_policy_acl" "http" {
  protocol = "tcp"
  ports = "80"
}

resource "enclave_policy" "testpolicy" {
  description = "this is a test"
  acl = [
    enclave_policy_acl.any,
    enclave_policy_acl.http
  ]
}

resource "enclave_dns_zone" "zone1" {
  name = "internal"
}

resource "enclave_dns_record" "record1"{
  name = "terraform-test"
  zone_id = enclave_dns_zone.zone1.id

}

resource "enclave_dns_zone" "zone2" {
  name = "external"
}

resource "enclave_dns_record" "record2"{
  name = "terraform-3"
  zone_id = enclave_dns_zone.zone2.id
}