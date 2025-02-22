/*
Copyright 2019 Nho Luong DevOps All rights reserved.

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

package addons

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"k8s.io/minikube/pkg/minikube/assets"
	"k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/minikube/localpath"
	"k8s.io/minikube/pkg/minikube/tests"
)

func createTestProfile(t *testing.T) string {
	t.Helper()
	td := t.TempDir()

	t.Setenv(localpath.MinikubeHome, td)

	// Not necessary, but it is a handy random alphanumeric
	name := filepath.Base(td)
	if err := os.MkdirAll(config.ProfileFolderPath(name), 0777); err != nil {
		t.Fatalf("error creating temporary directory")
	}

	cc := &config.ClusterConfig{
		Name:             name,
		CPUs:             2,
		Memory:           2500,
		KubernetesConfig: config.KubernetesConfig{},
		Nodes:            []config.Node{{ControlPlane: true}},
	}

	if err := config.DefaultLoader.WriteConfigToFile(name, cc); err != nil {
		t.Fatalf("error creating temporary profile config: %v", err)
	}
	return name
}

func TestIsAddonAlreadySet(t *testing.T) {
	cc := &config.ClusterConfig{
		Name:  "test",
		Nodes: []config.Node{{ControlPlane: true}},
	}

	if err := Set(cc, "registry", "true"); err != nil {
		t.Errorf("unable to set registry true: %v", err)
	}
	if !assets.Addons["registry"].IsEnabled(cc) {
		t.Errorf("expected registry to be enabled")
	}

	if assets.Addons["ingress"].IsEnabled(cc) {
		t.Errorf("expected ingress to not be enabled")
	}

}

func TestDisableUnknownAddon(t *testing.T) {
	cc := &config.ClusterConfig{
		Name:  "test",
		Nodes: []config.Node{{ControlPlane: true}},
	}

	if err := Set(cc, "InvalidAddon", "false"); err == nil {
		t.Fatalf("Disable did not return error for unknown addon")
	}
}

func TestEnableUnknownAddon(t *testing.T) {
	cc := &config.ClusterConfig{
		Name:  "test",
		Nodes: []config.Node{{ControlPlane: true}},
	}

	if err := Set(cc, "InvalidAddon", "true"); err == nil {
		t.Fatalf("Enable did not return error for unknown addon")
	}
}

func TestSetAndSave(t *testing.T) {
	profile := createTestProfile(t)

	// enable
	if err := SetAndSave(profile, "dashboard", "true"); err != nil {
		t.Errorf("Disable returned unexpected error: %v", err)
	}

	c, err := config.DefaultLoader.LoadConfigFromFile(profile)
	if err != nil {
		t.Errorf("unable to load profile: %v", err)
	}
	if c.Addons["dashboard"] != true {
		t.Errorf("expected dashboard to be enabled")
	}

	// disable
	if err := SetAndSave(profile, "dashboard", "false"); err != nil {
		t.Errorf("Disable returned unexpected error: %v", err)
	}

	c, err = config.DefaultLoader.LoadConfigFromFile(profile)
	if err != nil {
		t.Errorf("unable to load profile: %v", err)
	}
	if c.Addons["dashboard"] != false {
		t.Errorf("expected dashboard to be enabled")
	}
}

func TestStartWithAddonsEnabled(t *testing.T) {
	// this test will write a config.json into MinikubeHome, create a temp dir for it
	tests.MakeTempDir(t)

	cc := &config.ClusterConfig{
		Name:             "start",
		CPUs:             2,
		Memory:           2500,
		KubernetesConfig: config.KubernetesConfig{},
		Nodes:            []config.Node{{ControlPlane: true}},
	}

	toEnable := ToEnable(cc, map[string]bool{}, []string{"dashboard"})
	enabled := make(chan []string, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go Enable(&wg, cc, toEnable, enabled)
	wg.Wait()
	if ea, ok := <-enabled; ok {
		UpdateConfigToEnable(cc, ea)
	}

	if !assets.Addons["dashboard"].IsEnabled(cc) {
		t.Errorf("expected dashboard to be enabled")
	}
}

func TestStartWithAllAddonsDisabled(t *testing.T) {
	// this test will write a config.json into MinikubeHome, create a temp dir for it
	tests.MakeTempDir(t)

	cc := &config.ClusterConfig{
		Name:             "start",
		CPUs:             2,
		Memory:           2500,
		KubernetesConfig: config.KubernetesConfig{},
		Nodes:            []config.Node{{ControlPlane: true}},
	}

	UpdateConfigToDisable(cc)

	for name := range assets.Addons {
		if assets.Addons[name].IsEnabled(cc) {
			t.Errorf("expected %s to be disabled", name)
		}
	}
}
