# Copyright 2019 Nho Luong DevOps All rights reserved.
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

mkdir -p out

gsutil.cmd -m cp -r gs://minikube-builds/$env:MINIKUBE_LOCATION/common.ps1 out/

$driver="hyperv"
$timeout="180m"
$env:JOB_NAME="Hyper-V_Windows"
$env:EXTERNAL="yes"

. ./out/common.ps1
