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

package download

import (
	"fmt"
	"io/fs"
	"os"
	"sync"
	"testing"
	"time"

	"k8s.io/minikube/pkg/minikube/constants"
)

// Force download tests to run in serial.
func TestDownload(t *testing.T) {
	t.Run("BinaryDownloadPreventsMultipleDownload", testBinaryDownloadPreventsMultipleDownload)
	t.Run("PreloadDownloadPreventsMultipleDownload", testPreloadDownloadPreventsMultipleDownload)
	t.Run("ImageToCache", testImageToCache)
	t.Run("PreloadNotExists", testPreloadNotExists)
	t.Run("PreloadChecksumMismatch", testPreloadChecksumMismatch)
	t.Run("PreloadExistsCaching", testPreloadExistsCaching)
	t.Run("PreloadWithCachedSizeZero", testPreloadWithCachedSizeZero)
}

// Returns a mock function that sleeps before incrementing `downloadsCounter` and creates the requested file.
func mockSleepDownload(downloadsCounter *int) func(src, dst string) error {
	return func(src, dst string) error {
		// Sleep for 200ms to assure locking must have occurred.
		time.Sleep(time.Millisecond * 200)
		*downloadsCounter++
		return CreateDstDownloadMock(src, dst)
	}
}

func testBinaryDownloadPreventsMultipleDownload(t *testing.T) {
	downloadNum := 0
	DownloadMock = mockSleepDownload(&downloadNum)

	checkCache = func(_ string) (fs.FileInfo, error) {
		if downloadNum == 0 {
			return nil, fmt.Errorf("some error")
		}
		return nil, nil
	}

	var group sync.WaitGroup
	group.Add(2)
	dlCall := func() {
		if _, err := Binary("kubectl", "v1.20.2", "linux", "amd64", ""); err != nil {
			t.Errorf("Failed to download binary: %+v", err)
		}
		group.Done()
	}

	go dlCall()
	go dlCall()

	group.Wait()

	if downloadNum != 1 {
		t.Errorf("Expected only 1 download attempt but got %v!", downloadNum)
	}
}

func testPreloadDownloadPreventsMultipleDownload(t *testing.T) {
	downloadNum := 0
	DownloadMock = mockSleepDownload(&downloadNum)
	f, err := os.CreateTemp("", "preload")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err := f.Write([]byte("data")); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	checkCache = func(_ string) (fs.FileInfo, error) {
		if downloadNum == 0 {
			return nil, fmt.Errorf("some error")
		}
		return os.Stat(f.Name())
	}
	checkPreloadExists = func(_, _, _ string, _ ...bool) bool { return true }
	getChecksum = func(_, _ string) ([]byte, error) { return []byte("check"), nil }
	ensureChecksumValid = func(_, _, _ string, _ []byte) error { return nil }

	var group sync.WaitGroup
	group.Add(2)
	dlCall := func() {
		if err := Preload(constants.DefaultKubernetesVersion, constants.Docker, "docker"); err != nil {
			t.Logf("Failed to download preload: %+v (may be ok)", err)
		}
		group.Done()
	}

	go dlCall()
	go dlCall()

	group.Wait()

	if downloadNum != 1 {
		t.Errorf("Expected only 1 download attempt but got %v!", downloadNum)
	}
}

func testPreloadNotExists(t *testing.T) {
	downloadNum := 0
	DownloadMock = mockSleepDownload(&downloadNum)

	checkCache = func(_ string) (fs.FileInfo, error) { return nil, fmt.Errorf("cache not found") }
	checkPreloadExists = func(_, _, _ string, _ ...bool) bool { return false }
	getChecksum = func(_, _ string) ([]byte, error) { return []byte("check"), nil }
	ensureChecksumValid = func(_, _, _ string, _ []byte) error { return nil }

	err := Preload(constants.DefaultKubernetesVersion, constants.Docker, "docker")
	if err != nil {
		t.Errorf("Expected no error when preload exists")
	}

	if downloadNum != 0 {
		t.Errorf("Expected no download attempt but got %v!", downloadNum)
	}
}

func testPreloadChecksumMismatch(t *testing.T) {
	downloadNum := 0
	DownloadMock = mockSleepDownload(&downloadNum)

	checkCache = func(_ string) (fs.FileInfo, error) { return nil, fmt.Errorf("cache not found") }
	checkPreloadExists = func(_, _, _ string, _ ...bool) bool { return true }
	getChecksum = func(_, _ string) ([]byte, error) { return []byte("check"), nil }
	ensureChecksumValid = func(_, _, _ string, _ []byte) error {
		return fmt.Errorf("checksum mismatch")
	}

	err := Preload(constants.DefaultKubernetesVersion, constants.Docker, "docker")
	expectedErrMsg := "checksum mismatch"
	if err == nil {
		t.Errorf("Expected error when checksum mismatches")
	} else if err.Error() != expectedErrMsg {
		t.Errorf("Expected error to be %s, got %s", expectedErrMsg, err.Error())
	}
}

func testImageToCache(t *testing.T) {
	downloadNum := 0
	DownloadMock = mockSleepDownload(&downloadNum)

	checkImageExistsInCache = func(_ string) bool { return downloadNum > 0 }

	var group sync.WaitGroup
	group.Add(2)
	dlCall := func() {
		if err := ImageToCache("testimg"); err != nil {
			t.Errorf("Failed to download preload: %+v", err)
		}
		group.Done()
	}

	go dlCall()
	go dlCall()

	group.Wait()

	if downloadNum != 1 {
		t.Errorf("Expected only 1 download attempt but got %v!", downloadNum)
	}
}

// Validates that preload existence checks correctly caches values retrieved by remote checks.
func testPreloadExistsCaching(t *testing.T) {
	checkCache = func(_ string) (fs.FileInfo, error) {
		return nil, fmt.Errorf("cache not found")
	}
	doesPreloadExist := false
	checkCalled := false
	checkRemotePreloadExists = func(_, _ string) bool {
		checkCalled = true
		return doesPreloadExist
	}
	existence := PreloadExists("v1", "c1", "docker", true)
	if existence || !checkCalled {
		t.Errorf("Expected preload not to exist and a check to be performed. Existence: %v, Check: %v", existence, checkCalled)
	}
	checkCalled = false
	existence = PreloadExists("v1", "c1", "docker", true)
	if existence || checkCalled {
		t.Errorf("Expected preload not to exist and no check to be performed. Existence: %v, Check: %v", existence, checkCalled)
	}
	doesPreloadExist = true
	checkCalled = false
	existence = PreloadExists("v2", "c1", "docker", true)
	if !existence || !checkCalled {
		t.Errorf("Expected preload to exist and a check to be performed. Existence: %v, Check: %v", existence, checkCalled)
	}
	checkCalled = false
	existence = PreloadExists("v2", "c2", "docker", true)
	if !existence || !checkCalled {
		t.Errorf("Expected preload to exist and a check to be performed. Existence: %v, Check: %v", existence, checkCalled)
	}
	checkCalled = false
	existence = PreloadExists("v2", "c2", "docker", true)
	if !existence || checkCalled {
		t.Errorf("Expected preload to exist and no check to be performed. Existence: %v, Check: %v", existence, checkCalled)
	}
}

func testPreloadWithCachedSizeZero(t *testing.T) {
	downloadNum := 0
	DownloadMock = mockSleepDownload(&downloadNum)
	f, err := os.CreateTemp("", "preload")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	checkCache = func(_ string) (fs.FileInfo, error) { return os.Stat(f.Name()) }
	checkPreloadExists = func(_, _, _ string, _ ...bool) bool { return true }
	getChecksum = func(_, _ string) ([]byte, error) { return []byte("check"), nil }
	ensureChecksumValid = func(_, _, _ string, _ []byte) error { return nil }

	if err := Preload(constants.DefaultKubernetesVersion, constants.Docker, "docker"); err != nil {
		t.Errorf("Expected no error with cached preload of size zero")
	}

	if downloadNum != 1 {
		t.Errorf("Expected only 1 download attempt but got %v!", downloadNum)
	}
}
