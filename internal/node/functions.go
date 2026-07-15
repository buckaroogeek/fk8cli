/*
Copyright 2025 Bradley G Smith >bradley.g.smith@gmail.com>
SPDX-License-Identifier: MIT

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Kubernetes installation functions
package node

import (
	"fmt"
	"os"
	"os/user"

	"buckaroogeek.com/fk8cli/internal/initialize"
	"github.com/bitfield/script"
)

// echo basic configuration info to log file
func echoConfig(cfg *initialize.Config) error {

	// Append to log file
	f, err := os.OpenFile(cfg.FileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// configuration
	_, err = fmt.Fprintf(f, "Target configuration:\n"+
		"  Version: %s",
		cfg.Tag())

	_, err = fmt.Fprintf(f, "\n\nInstallation Steps\n")

	if err != nil {
		return err
	}
	return nil
}

// Process an install/configuration step
func process(cfg *initialize.Config, fn func(cfg *initialize.Config) error, msg string) error {
	// msg to stdout
	fmt.Println(" ..." + msg)

	// msg to log file
	// open log file
	f, err := os.OpenFile(cfg.FileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(f, "\n*** %s ***\n", msg); err != nil {
		return err
	}
	f.Close()

	// execute process function
	return fn(cfg)
}

// execute function with sudo
func sudoexec(cmd string, cfg *initialize.Config) error {
	return exec("sudo "+cmd, cfg)
}

// execute function - adapted from
// https://raw.githubusercontent.com/ccollicutt/go-install-kubernetes/refs/heads/main/pkg/exec/exec.go
func exec(cmd string, cfg *initialize.Config) error {

	// open log file
	f, err := os.OpenFile(cfg.FileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Print cmd to stdout
	fmt.Printf("$ %s\n", cmd)

	// Append command to log file
	if _, err := fmt.Fprintf(f, "\n$ %s\n", cmd); err != nil {
		return err
	}

	// If verbose, also print to stdout

	// execute command for this step
	output, err := script.Exec(cmd).String()
	if err != nil {
		return err
	}

	// Append output to log file
	if _, err := fmt.Fprintf(f, "%s\n", output); err != nil {
		return err
	}

	// If verbose, also print output to stdout
	if cfg.IsVerbose() {
		fmt.Printf("\n%s\n", output)
	}

	return nil
}

// install rpms
func installPackages(cfg *initialize.Config) error {
	cmd := fmt.Sprintf("dnf install -y %s", cfg.Rpms()+cfg.Rpmfiles())
	return sudoexec(cmd, cfg)
}

// configure CRI (cri-o by default)
func configureCRI(cfg *initialize.Config) error {
	// load and persist br_netfilter (overlay no longer in k8s docs)
	content := "br_netfilter\n"
	err := idempotentWrite("/etc/modules-load.d/k8s.conf", []byte(content), cfg)
	if err != nil {
		return err
	}

	// load br_netfilter"
	cmd := "modprobe br_netfilter"
	err = sudoexec(cmd, cfg)
	if err != nil {
		return err
	}

	// ipv4 packet forwarding
	content = "net.ipv4.ip_forward = 1\n"
	err = idempotentWrite("/etc/sysctl.d/99-k8s.conf", []byte(content), cfg)
	if err != nil {
		return err
	}

	// apply sysctl parameters
	cmd = "sysctl --system"
	err = sudoexec(cmd, cfg)
	if err != nil {
		return err
	}

	// start and enable CRI
	cmd = "systemctl enable --now crio"
	return sudoexec(cmd, cfg)
}

// configure kubernetes via kubeadm
func configureK8S(cfg *initialize.Config) error {

	// enable kubelet, kubeadm will start as needed
	cmd := "systemctl enable kubelet"
	err := sudoexec(cmd, cfg)
	if err != nil {
		return err
	}

	// execute kubeadm
	cmd = "kubeadm init --pod-network-cidr=10.244.0.0/16"
	err = sudoexec(cmd, cfg)
	if err != nil {
		return err
	}

	// configure kubectl access
	cmd = "mkdir -p $HOME/.kube"
	err = exec(cmd, cfg)
	if err != nil {
		return err
	}
	cmd = "cp  /etc/kubernetes/admin.conf $HOME/.kube/config"
	err = sudoexec(cmd, cfg)
	if err != nil {
		return err
	}

	// use the user name to pull Uid and Gid
	sudouser, err := user.Lookup(cfg.User())
	if err != nil {
		return err
	}
	cmd = "chown " + sudouser.Uid + ":" + sudouser.Gid + " $HOME/.kube/config"
	err = sudoexec(cmd, cfg)
	if err != nil {
		return err
	}

	// set taint if toggled
	if cfg.GetTaint() {
		cmd = "kubectl taint nodes --all node-role.kubernetes.io/control-plane-"
		err = exec(cmd, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

// configure the network interface (flannel by default)
func configureCNI(cfg *initialize.Config) error {
	// configure flannel
	cmd := "kubectl apply -f https://github.com/coreos/flannel/raw/master/Documentation/kube-flannel.yml"
	return exec(cmd, cfg)
}
