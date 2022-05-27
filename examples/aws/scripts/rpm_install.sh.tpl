#!/bin/bash
dnf -y install dnf-plugins-core
dnf config-manager --add-repo https://packages.enclave.io/rpm/enclave.repo

dnf install enclave -y --refresh

sudo enclave enrol ${enclave_key}