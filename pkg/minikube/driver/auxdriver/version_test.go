/*
Copyright 2020 Nho Luong DevOps All rights reserved.

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

package auxdriver

import (
	"testing"

	"github.com/blang/semver/v4"
	"k8s.io/minikube/pkg/minikube/driver"
)

func Test_minDriverVersion(t *testing.T) {

	tests := []struct {
		desc   string
		driver string
		mkV    string
		want   semver.Version
	}{
		{"Hyperkit", driver.HyperKit, "1.1.1", *minHyperkitVersion},
		{"Invalid", "_invalid_", "1.1.1", v("1.1.1")},
		{"KVM2", driver.KVM2, "1.1.1", v("1.1.1")},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got := minAcceptableDriverVersion(tt.driver, v(tt.mkV)); !got.EQ(tt.want) {
				t.Errorf("Invalid min supported version, got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func v(s string) semver.Version {
	r, err := semver.New(s)
	if err != nil {
		panic(err)
	}
	return *r
}
