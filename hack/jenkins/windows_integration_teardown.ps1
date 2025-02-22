# Copyright 2021 Nho Luong DevOps All rights reserved.
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

$test_home="$env:HOMEDRIVE$env:HOMEPATH\minikube-integration"

if ($driver -eq "docker") {
  # Remove unused images and containers
  docker system prune --all --force --volumes

  # Just shutdown Docker, it's safer than anything else
  Get-Process "*Docker Desktop*" | Stop-Process
}

rm -r -Force $test_home
C:\jenkins\windows_cleanup_and_reboot.ps1
