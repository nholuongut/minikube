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

set -e

for pkg in $(go list ./pkg/minikube/...); do
    go test -v $pkg --coverprofile=out/coverage.out
    go tool cover -html out/coverage.out
done