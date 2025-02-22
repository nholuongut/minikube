/*
Copyright 2016 Nho Luong DevOps All rights reserved.

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

package driver

import (
	"os/exec"
)

// supportedDrivers is a list of supported drivers on Linux.
var supportedDrivers = []string{
	VirtualBox,
	KVM2,
	QEMU2,
	QEMU,
	VMware,
	None,
	Docker,
	Podman,
	SSH,
}

// VBoxManagePath returns the path to the VBoxManage command
func VBoxManagePath() string {
	cmd := "VBoxManage"
	if path, err := exec.LookPath(cmd); err == nil {
		return path
	}
	return cmd
}
