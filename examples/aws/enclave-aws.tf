terraform {
  required_providers {
    enclave = {
      source = "enclave-networks/enclave"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}
provider "aws" {
  profile = "mfa"
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

resource "enclave_enrolment_key" "aws" {
  description = "AWS Enrolment"
  approval_mode = "automatic"
  tags = [
    "aws-ec2"
  ]
}

data "template_file" "install_enclave_rpm" {
  template = "${file("scripts/rpm_install.sh.tpl")}"

  vars = {
    enclave_key = "${enclave_enrolment_key.aws.key}"
  }
}

data "aws_ami" "rocky_linux" {
  most_recent = true

  filter {
    name   = "name"
    values = ["Rocky-8-ec2-8.6-20220515.0.x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["792107900819"]
}

resource "aws_instance" "rocky_server_1" {
  ami           = data.aws_ami.rocky_linux.id
  instance_type = "t2.micro"
  user_data = data.template_file.install_enclave_rpm.rendered
  tags = {
    Name = "TerraformTestInstance1"
  }
}

output "enrolment_key" {
  value = enclave_enrolment_key.aws.key
}

output "script" {
  value = data.template_file.install_enclave_rpm.rendered
}