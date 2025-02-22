#!/bin/bash

# Copyright 2016 Nho Luong DevOps All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# This script runs the integration tests on a Linux machine for the KVM Driver

# The script expects the following env variables:
# MINIKUBE_LOCATION: GIT_COMMIT from upstream build.
# COMMIT: Actual commit ID from upstream build
# EXTRA_BUILD_ARGS (optional): Extra args to be passed into the minikube integrations tests
# access_token: The GitHub API access token. Injected by the Jenkins credential provider.

set -e

OS="linux"
ARCH="amd64"
DRIVER="kvm2"
JOB_NAME="KVM_Linux"

# We pick kvm as our gvisor testbed because it is fast & reliable
EXTRA_TEST_ARGS="-gvisor"

sudo apt-get update
sudo apt-get -y install qemu-system libvirt-clients libvirt-daemon-system ebtables iptables dnsmasq
sudo adduser jenkins libvirt || true

source ./common.sh
