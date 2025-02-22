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

package bsutil

import (
	"strings"

	"github.com/blang/semver/v4"
	"k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/util"
)

// versionIsBetween checks if a version is between (or including) two given versions
func versionIsBetween(version, gte, lte semver.Version) bool {
	if gte.NE(semver.Version{}) && !version.GTE(gte) {
		return false
	}
	if lte.NE(semver.Version{}) && !version.LTE(lte) {
		return false
	}

	return true
}

var versionSpecificOpts = []config.VersionedExtraOption{
	config.NewUnversionedOption(Kubelet, "bootstrap-kubeconfig", "/etc/kubernetes/bootstrap-kubelet.conf"),
	config.NewUnversionedOption(Kubelet, "config", "/var/lib/kubelet/config.yaml"),
	config.NewUnversionedOption(Kubelet, "kubeconfig", "/etc/kubernetes/kubelet.conf"),
	{
		Option: config.ExtraOption{
			Component: Apiserver,
			Key:       "enable-admission-plugins",
			Value:     strings.Join(util.DefaultAdmissionControllers, ","),
		},
		GreaterThanOrEqual: semver.MustParse("1.14.0-alpha.0"),
	},
	{
		Option: config.ExtraOption{
			Component: ControllerManager,
			Key:       "allocate-node-cidrs",
			Value:     "true",
		},
		GreaterThanOrEqual: semver.MustParse("1.14.0"),
	},
	{
		Option: config.ExtraOption{
			Component: ControllerManager,
			Key:       "leader-elect",
			Value:     "false",
		},
		GreaterThanOrEqual: semver.MustParse("1.14.0"),
	},
	{
		Option: config.ExtraOption{
			Component: Scheduler,
			Key:       "leader-elect",
			Value:     "false",
		},
		GreaterThanOrEqual: semver.MustParse("1.14.0"),
	},
}
