---
title: "Addons"
weight: 4
description: >
  How to develop minikube addons
---

## Creating a new addon

To create an addon, first fork the minikube repository, and check out your fork:

```shell
git clone git@github.com:<username>/minikube.git
```

Then go into the source directory:

```shell
cd minikube
```

Create a subdirectory:

```shell
mkdir deploy/addons/<addon name>
```

Add your manifest YAML's to the directory you have created:

```shell
cp *.yaml deploy/addons/<addon name>
```

Note: If the addon never needs authentication to GCP, then consider adding the following label to the pod's yaml:

```yaml
gcp-auth-skip-secret: "true"
```

To make the addon appear in `minikube addons list`, add it to `pkg/addons/config.go`. Here is the entry used by the `registry` addon, which will work for any addon which does not require custom code:

```go
  {
    name:      "registry",
    set:       SetBool,
    callbacks: []setFn{EnableOrDisableAddon},
  },
```

Next, add all required files using `//go:embed` directives to a new embed.FS variable in `deploy/addons/assets.go`. Here is the entry used by the `csi-hostpath-driver` addon:

```go
	// CsiHostpathDriverAssets assets for csi-hostpath-driver addon
	//go:embed csi-hostpath-driver/deploy/*.tmpl csi-hostpath-driver/rbac/*.tmpl
	CsiHostpathDriverAssets embed.FS
```

Then, add into `pkg/minikube/assets/addons.go` the list of files to copy into the cluster, including manifests. Here is the entry used by the `registry` addon:

```go
  "registry": NewAddon([]*BinAsset{
    MustBinAsset(addons.RegistryAssets,
      "registry/registry-rc.yaml.tmpl",
      vmpath.GuestAddonsDir,
      "registry-rc.yaml",
      "0640",
      false),
    MustBinAsset(addons.RegistryAssets,
      "registry/registry-svc.yaml.tmpl",
      vmpath.GuestAddonsDir,
      "registry-svc.yaml",
      "0640",
      false),
    MustBinAsset(addons.RegistryAssets,
      "registry/registry-proxy.yaml.tmpl",
      vmpath.GuestAddonsDir,
      "registry-proxy.yaml",
      "0640",
      false),
  }, false, "registry", "google"),
```

The `MustBinAsset` arguments are:

* asset variable (typically present in `deploy/addons/assets.go`)
* source filename
* destination directory (typically `vmpath.GuestAddonsDir`)
* destination filename
* permissions (typically `0640`)
* boolean value representing if template substitution is required (often `false`)

The boolean value on the last line is whether the addon should be enabled by default. This should always be `false`. In addition, following the addon name on the last line is the maintainer field. This is meant to inform users about the controlling party of an addon's images. In the case above, the maintainer is Google, since the registry addon uses images that Google controls. When creating a new addon, the source of the images should be contacted and requested whether they are willing to be the point of contact for this addon before being put. If the source does not accept the responsibility, leaving the maintainer field empty is acceptable.

To see other examples, see the [addons commit history](https://github.com/nholuongut/minikube/commits/master/deploy/addons) for other recent examples.

## "addons open" support

If your addon contains a NodePort Service, please add the `kubernetes.io/minikube-addons-endpoint: <addon name>` label, which is used by the  `minikube addons open` command:

```yaml
apiVersion: v1
kind: Service
metadata:
 labels:
    kubernetes.io/minikube-addons-endpoint: <addon name>
```

NOTE: `minikube addons open` currently only works for the `kube-system` namespace: [#8089](https://github.com/nholuongut/minikube/issues/8089).

## Testing addon changes

Rebuild the minikube binary and apply the addon with extra logging enabled:

```shell
make && make test && ./out/minikube addons enable <addon name> --alsologtostderr
```

Please note that you must run `make` each time you change your YAML files. To disable the addon when new changes are made, run:

```shell
./out/minikube addons disable <addon name> --alsologtostderr
```

## Sending out your PR

Once you have tested your addon, click on [new pull request](https://github.com/nholuongut/minikube/compare) to send us your PR!
