/*
Copyright 2025 Bradley G Smith >bradley.g.smith@gmail.com>
SPDX-License-Identifier: MIT

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Kubernetes installation controller
// Theory of operation
// 1. disable firewall (tbd add firewall support)
// 2. install CRI (cri-o default)
// 2a. IPv4 packet forwarding (ipv6 tbd)
// 2b. bridge filter
// 2c. enable cri-o service
// 3.  install kubernetes rpms
// 3a. enable kubelet service
// 4. kubeadm init
// 5 control plane config
// 5a. taint (only if single node or -t flag set)
// 5b. CNI pod network add on (flannel by default)
package node

import (
	"fmt"

	"buckaroogeek.com/fk8cli/internal/initialize"
)

// Manage installation process
func Install(cfg *initialize.Config) error {
	fmt.Println("\nStarting install for K8S version: ", cfg.Version())

	//
	echoConfig(cfg)

	// install packages
	err := process(cfg, installPackages, "Install all packages")
	if err != nil {
		return err
	}

	// configure CRI (cri-o by default)
	err = process(cfg, configureCRI, "Configure CRI")
	if err != nil {
		return err
	}

	// configure Kubernetes
	err = process(cfg, configureK8S, "Configure Kubernetes")
	if err != nil {
		return err
	}

	// configure CNI
	err = process(cfg, configureCNI, "Configure CNI")
	if err != nil {
		return err
	}

	// verify functioning
	cmd := "kubectl get pods --all-namespaces"
	err = exec(cmd, cfg)
	if err != nil {
		return err
	}

	// k8s install done
	return nil
}
