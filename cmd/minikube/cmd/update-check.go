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

package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/minikube/pkg/minikube/exit"
	"k8s.io/minikube/pkg/minikube/notify"
	"k8s.io/minikube/pkg/minikube/out"
	"k8s.io/minikube/pkg/minikube/reason"
	"k8s.io/minikube/pkg/version"
)

var updateCheckCmd = &cobra.Command{
	Use:   "update-check",
	Short: "Print current and latest version number",
	Long:  `Print current and latest version number`,
	Run: func(_ *cobra.Command, _ []string) {
		url := notify.GithubMinikubeReleasesURL
		r, err := notify.AllVersionsFromURL(url)
		if err != nil {
			exit.Error(reason.InetVersionUnavailable, "Unable to fetch latest version info", err)
		}

		if len(r.Releases) < 1 {
			exit.Message(reason.InetVersionEmpty, "Update server returned an empty list")
		}

		out.Ln("CurrentVersion: %s", version.GetVersion())
		out.Ln("LatestVersion: %s", r.Releases[0].Name)
	},
}
