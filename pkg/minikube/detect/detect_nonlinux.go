//go:build !linux

/*
Copyright 2022 Nho Luong DevOps All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package detect

// cgroupVersion returns cgroups v1 for non-linux OS host machine (where minikube runs).
func cgroupVersion() string {
	return "v1"
}

// Assume 9p is supported by non linux apps
func IsNinePSupported() bool {
	return true
}
