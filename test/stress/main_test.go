/*
Copyright 2021 Nho Luong DevOps All rights reserved.

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

package stress

import (
	"flag"
	"os"
	"strings"
	"testing"
)

var startArgs = flag.String("start-args", "", "Arguments to pass to minikube start")
var upgradeFrom = flag.String("upgrade-from", "v1.11.0", "The version of minikube to start with, and upgrade from.")
var loops = flag.Int("loops", 20, "The number of times to run the test")

func TestMain(m *testing.M) {
	flag.Parse()
	if !strings.HasPrefix(*upgradeFrom, "v") {
		*upgradeFrom = "v" + *upgradeFrom
	}
	os.Exit(m.Run())
}
