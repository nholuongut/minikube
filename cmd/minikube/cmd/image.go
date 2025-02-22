/*
Copyright 2017 Nho Luong DevOps All rights reserved.

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
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/minikube/exit"
	"k8s.io/minikube/pkg/minikube/image"
	"k8s.io/minikube/pkg/minikube/machine"
	"k8s.io/minikube/pkg/minikube/out"
	"k8s.io/minikube/pkg/minikube/reason"
	docker "k8s.io/minikube/third_party/go-dockerclient"
)

var (
	allNodes bool
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image COMMAND",
	Short: "Manage images",
}

var (
	pull       bool
	imgDaemon  bool
	imgRemote  bool
	overwrite  bool
	tag        string
	push       bool
	dockerFile string
	buildEnv   []string
	buildOpt   []string
	format     string
)

func saveFile(r io.Reader) (string, error) {
	tmp, err := os.CreateTemp("", "build.*.tar")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(tmp, r)
	if err != nil {
		return "", err
	}
	err = tmp.Close()
	if err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

// loadImageCmd represents the image load command
var loadImageCmd = &cobra.Command{
	Use:     "load IMAGE | ARCHIVE | -",
	Short:   "Load an image into minikube",
	Long:    "Load an image into minikube",
	Example: "minikube image load image\nminikube image load image.tar",
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			exit.Message(reason.Usage, "Please provide an image in your local daemon to load into minikube via <minikube image load IMAGE_NAME>")
		}
		// Cache and load images into container runtime
		profile, err := config.LoadProfile(viper.GetString(config.ProfileName))
		if err != nil {
			exit.Error(reason.Usage, "loading profile", err)
		}

		if pull {
			// Pull image from remote registry, without doing any caching except in container runtime.
			// This is similar to daemon.Image but it is done by the container runtime in the cluster.
			if err := machine.PullImages(args, profile); err != nil {
				exit.Error(reason.GuestImageLoad, "Failed to pull image", err)
			}
			return
		}

		var local bool
		if imgRemote || imgDaemon {
			local = false
		} else {
			for _, img := range args {
				if img == "-" { // stdin
					local = true
					imgDaemon = false
					imgRemote = false
				} else if strings.HasPrefix(img, "/") || strings.HasPrefix(img, ".") {
					local = true
					imgDaemon = false
					imgRemote = false
				} else if _, err := os.Stat(img); err == nil {
					local = true
					imgDaemon = false
					imgRemote = false
				}
			}

			if !local {
				imgDaemon = true
				imgRemote = true
			}
		}

		if args[0] == "-" {
			tmp, err := saveFile(os.Stdin)
			if err != nil {
				exit.Error(reason.GuestImageLoad, "Failed to save stdin", err)
			}
			args = []string{tmp}
		}

		if imgDaemon || imgRemote {
			image.UseDaemon(imgDaemon)
			image.UseRemote(imgRemote)
			if err := machine.CacheAndLoadImages(args, []*config.Profile{profile}, overwrite); err != nil {
				exit.Error(reason.GuestImageLoad, "Failed to load image", err)
			}
		} else if local {
			// Load images from local files, without doing any caching or checks in container runtime
			// This is similar to tarball.Image but it is done by the container runtime in the cluster.
			if err := machine.DoLoadImages(args, []*config.Profile{profile}, "", overwrite); err != nil {
				exit.Error(reason.GuestImageLoad, "Failed to load image", err)
			}
		}
	},
}

func readFile(w io.Writer, tmp string) error {
	r, err := os.Open(tmp)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	err = r.Close()
	if err != nil {
		return err
	}
	return nil
}

// saveImageCmd represents the image load command
var saveImageCmd = &cobra.Command{
	Use:     "save IMAGE [ARCHIVE | -]",
	Short:   "Save a image from minikube",
	Long:    "Save a image from minikube",
	Example: "minikube image save image\nminikube image save image image.tar",
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			exit.Message(reason.Usage, "Please provide an image in the container runtime to save from minikube via <minikube image save IMAGE_NAME>")
		}
		// Save images from container runtime
		profile, err := config.LoadProfile(viper.GetString(config.ProfileName))
		if err != nil {
			exit.Error(reason.Usage, "loading profile", err)
		}

		if len(args) > 1 {
			output = args[1]

			if args[1] == "-" {
				tmp, err := os.CreateTemp("", "image.*.tar")
				if err != nil {
					exit.Error(reason.GuestImageSave, "Failed to get temp", err)
				}
				tmp.Close()
				output = tmp.Name()
			}

			if err := machine.DoSaveImages([]string{args[0]}, output, []*config.Profile{profile}, ""); err != nil {
				exit.Error(reason.GuestImageSave, "Failed to save image", err)
			}

			if args[1] == "-" {
				err := readFile(os.Stdout, output)
				if err != nil {
					exit.Error(reason.GuestImageSave, "Failed to read temp", err)
				}
				os.Remove(output)
			}
		} else {
			if err := machine.SaveAndCacheImages([]string{args[0]}, []*config.Profile{profile}); err != nil {
				exit.Error(reason.GuestImageSave, "Failed to save image", err)
			}
			if imgDaemon || imgRemote {
				image.UseDaemon(imgDaemon)
				image.UseRemote(imgRemote)
				err := image.UploadCachedImage(args[0])
				if err != nil {
					exit.Error(reason.GuestImageSave, "Failed to save image", err)
				}
			}
		}
	},
}

var removeImageCmd = &cobra.Command{
	Use:   "rm IMAGE [IMAGE...]",
	Short: "Remove one or more images",
	Example: `
$ minikube image rm image busybox

$ minikube image unload image busybox
`,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"remove", "unload"},
	Run: func(_ *cobra.Command, args []string) {
		profile, err := config.LoadProfile(viper.GetString(config.ProfileName))
		if err != nil {
			exit.Error(reason.Usage, "loading profile", err)
		}
		if err := machine.RemoveImages(args, profile); err != nil {
			exit.Error(reason.GuestImageRemove, "Failed to remove image", err)
		}
	},
}

var pullImageCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull images",
	Example: `
$ minikube image pull busybox
`,
	Run: func(_ *cobra.Command, args []string) {
		profile, err := config.LoadProfile(viper.GetString(config.ProfileName))
		if err != nil {
			exit.Error(reason.Usage, "loading profile", err)
		}

		if err := machine.PullImages(args, profile); err != nil {
			exit.Error(reason.GuestImagePull, "Failed to pull images", err)
		}
	},
}

func createTar(dir string) (string, error) {
	tar, err := docker.CreateTarStream(dir, dockerFile)
	if err != nil {
		return "", err
	}
	return saveFile(tar)
}

// buildImageCmd represents the image build command
var buildImageCmd = &cobra.Command{
	Use:     "build PATH | URL | -",
	Short:   "Build a container image in minikube",
	Long:    "Build a container image, using the container runtime.",
	Example: `minikube image build .`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) < 1 {
			exit.Message(reason.Usage, "Please provide a path or url to build")
		}
		// Build images into container runtime
		profile, err := config.LoadProfile(viper.GetString(config.ProfileName))
		if err != nil {
			exit.Error(reason.Usage, "loading profile", err)
		}

		img := args[0]
		var tmp string
		if img == "-" {
			tmp, err = saveFile(os.Stdin)
			if err != nil {
				exit.Error(reason.GuestImageBuild, "Failed to save stdin", err)
			}
			img = tmp
		} else {
			// If it is an URL, pass it as-is
			u, err := url.Parse(img)
			local := err == nil && u.Scheme == "" && u.Host == ""
			if runtime.GOOS == "windows" && filepath.VolumeName(img) != "" {
				local = true
			}
			if local {
				// If it's a directory, tar it
				info, err := os.Stat(img)
				if err == nil && info.IsDir() {
					tmp, err := createTar(img)
					if err != nil {
						exit.Error(reason.GuestImageBuild, "Failed to save dir", err)
					}
					img = tmp
				}
				// Otherwise, assume it's a tar
			}
		}
		if runtime.GOOS == "windows" && strings.Contains(dockerFile, "\\") {
			// if dockerFile is a DOS path, translate it into UNIX path
			// because we are going to build this image in UNIX environment
			out.Stringf("minikube detects that you are using DOS-style path %s. minikube will convert it to UNIX-style by replacing all \\ to /\n", dockerFile)
			dockerFile = strings.ReplaceAll(dockerFile, "\\", "/")
		}
		if err := machine.BuildImage(img, dockerFile, tag, push, buildEnv, buildOpt, []*config.Profile{profile}, allNodes, nodeName); err != nil {
			exit.Error(reason.GuestImageBuild, "Failed to build image", err)
		}
		if tmp != "" {
			os.Remove(tmp)
		}
	},
}

var listImageCmd = &cobra.Command{
	Use:   "ls",
	Short: "List images",
	Example: `
$ minikube image ls
`,
	Aliases: []string{"list"},
	Run: func(_ *cobra.Command, _ []string) {
		profile, err := config.LoadProfile(viper.GetString(config.ProfileName))
		if err != nil {
			exit.Error(reason.Usage, "loading profile", err)
		}

		if err := machine.ListImages(profile, format); err != nil {
			exit.Error(reason.GuestImageList, "Failed to list images", err)
		}
	},
}

var tagImageCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag images",
	Example: `
$ minikube image tag source target
`,
	Aliases: []string{"list"},
	Run: func(_ *cobra.Command, args []string) {
		if len(args) != 2 {
			exit.Message(reason.Usage, "Please provide source and target image")
		}
		profile, err := config.LoadProfile(viper.GetString(config.ProfileName))
		if err != nil {
			exit.Error(reason.Usage, "loading profile", err)
		}

		if err := machine.TagImage(profile, args[0], args[1]); err != nil {
			exit.Error(reason.GuestImageTag, "Failed to tag images", err)
		}
	},
}

var pushImageCmd = &cobra.Command{
	Use:   "push",
	Short: "Push images",
	Example: `
$ minikube image push busybox
`,
	Run: func(_ *cobra.Command, args []string) {
		profile, err := config.LoadProfile(viper.GetString(config.ProfileName))
		if err != nil {
			exit.Error(reason.Usage, "loading profile", err)
		}

		if err := machine.PushImages(args, profile); err != nil {
			exit.Error(reason.GuestImagePush, "Failed to push images", err)
		}
	},
}

func init() {
	loadImageCmd.Flags().BoolVar(&pull, "pull", false, "Pull the remote image (no caching)")
	loadImageCmd.Flags().BoolVar(&imgDaemon, "daemon", false, "Cache image from docker daemon")
	loadImageCmd.Flags().BoolVar(&imgRemote, "remote", false, "Cache image from remote registry")
	loadImageCmd.Flags().BoolVar(&overwrite, "overwrite", true, "Overwrite image even if same image:tag name exists")
	imageCmd.AddCommand(loadImageCmd)
	imageCmd.AddCommand(removeImageCmd)
	imageCmd.AddCommand(pullImageCmd)
	buildImageCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag to apply to the new image (optional)")
	buildImageCmd.Flags().BoolVar(&push, "push", false, "Push the new image (requires tag)")
	buildImageCmd.Flags().StringVarP(&dockerFile, "file", "f", "", "Path to the Dockerfile to use (optional)")
	buildImageCmd.Flags().StringArrayVar(&buildEnv, "build-env", nil, "Environment variables to pass to the build. (format: key=value)")
	buildImageCmd.Flags().StringArrayVar(&buildOpt, "build-opt", nil, "Specify arbitrary flags to pass to the build. (format: key=value)")
	buildImageCmd.Flags().StringVarP(&nodeName, "node", "n", "", "The node to build on. Defaults to the primary control plane.")
	buildImageCmd.Flags().BoolVar(&allNodes, "all", false, "Build image on all nodes.")
	imageCmd.AddCommand(buildImageCmd)
	saveImageCmd.Flags().BoolVar(&imgDaemon, "daemon", false, "Cache image to docker daemon")
	saveImageCmd.Flags().BoolVar(&imgRemote, "remote", false, "Cache image to remote registry")
	imageCmd.AddCommand(saveImageCmd)
	listImageCmd.Flags().StringVar(&format, "format", "short", "Format output. One of: short|table|json|yaml")
	imageCmd.AddCommand(listImageCmd)
	imageCmd.AddCommand(tagImageCmd)
	imageCmd.AddCommand(pushImageCmd)
}
