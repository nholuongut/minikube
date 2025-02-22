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

package driver

import (
	"fmt"
	"net"

	"k8s.io/klog/v2"
	"k8s.io/minikube/pkg/drivers/kic/oci"
	"k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/minikube/constants"
	"k8s.io/minikube/pkg/network"
)

// ControlPlaneEndpoint returns the location where callers can reach this cluster.
func ControlPlaneEndpoint(cc *config.ClusterConfig, cp *config.Node, driverName string) (string, net.IP, int, error) {
	if NeedsPortForward(driverName) {
		port, err := oci.ForwardedPort(cc.Driver, cc.Name, cp.Port)
		if err != nil {
			klog.Warningf("failed to get forwarded control plane port %v", err)
		}

		hostname := oci.DaemonHost(driverName)

		ips, err := net.LookupIP(hostname)
		if err != nil || len(ips) == 0 {
			return hostname, nil, port, fmt.Errorf("failed to lookup ip for %q", hostname)
		}

		// https://github.com/nholuongut/minikube/issues/3878
		if cc.KubernetesConfig.APIServerName != constants.APIServerName {
			hostname = cc.KubernetesConfig.APIServerName
		}

		return hostname, ips[0], port, nil
	}

	if IsQEMU(driverName) && network.IsBuiltinQEMU(cc.Network) {
		return "localhost", net.IPv4(127, 0, 0, 1), cc.APIServerPort, nil
	}

	// https://github.com/nholuongut/minikube/issues/3878
	hostname := cp.IP
	if cc.KubernetesConfig.APIServerName != constants.APIServerName {
		hostname = cc.KubernetesConfig.APIServerName
	}
	ips, err := net.LookupIP(cp.IP)
	if err != nil || len(ips) == 0 {
		return hostname, nil, cp.Port, fmt.Errorf("failed to lookup ip for %q", cp.IP)
	}
	return hostname, ips[0], cp.Port, nil
}

// AutoPauseProxyEndpoint returns the endpoint for the auto-pause (reverse proxy to api-sever)
func AutoPauseProxyEndpoint(cc *config.ClusterConfig, cp *config.Node, driverName string) (string, net.IP, int, error) {
	cp.Port = constants.AutoPauseProxyPort
	return ControlPlaneEndpoint(cc, cp, driverName)
}
