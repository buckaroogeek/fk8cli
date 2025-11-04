/*
Copyright 2025 Bradley G Smith >bradley.g.smith@gmail.com>
SPDX-License-Identifier: MIT

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Initialize cli
// Returns pointer to struct populated by flags

package initialize

import (
	"flag"
	"fmt"
	"os"

	"github.com/bitfield/script"
)

/* Sanity checks, parse flags and arguements */
func ParseFlags() *Config {
	cfg := NewConfig()

	flag.BoolVar(&cfg.isControl, "c", cfg.isControl, "Configure as a control plane node")
	flag.BoolVar(&cfg.isWorker, "w", cfg.isWorker, "Configure as a worker node")
	flag.BoolVar(&cfg.isSingle, "s", cfg.isSingle, "Configure as a single node (control plane + worker)")
	flag.BoolVar(&cfg.isQuiet, "q", cfg.isQuiet, "Enable verbose output")
	flag.BoolVar(&cfg.isGo, "y", cfg.isGo, "Proceed with installation")

	flag.Usage = showHelp
	flag.Parse()

	// check root access
	checkRoot()

	// user must not combine -s with -c or -w
	if cfg.isSingle {
		if cfg.isControl || cfg.isWorker {
			showHelpAndExit("Cannot combine -s with -c and/or -w", 1)
		}
		// toggle control and worker to true given isSingle is set
		cfg.isControl = true
		cfg.isWorker = true
	}

	// check that at least one role is set
	if !(cfg.isControl || cfg.isWorker) {
		showHelpAndExit("At least one role (control or worker) must be set", 1)
	}

	// Check arguments - only 1 allowed (version)
	checkArgs()

	// Check for version argument
	checkVersion(cfg)

	// Build list of rpms to install
	buildRPMList(cfg)

	// Show configuration
	showConfiguration(cfg)

	// if dryrun show configuration and exit
	if !cfg.isGo {
		showHelpAndExit("Dry run", 0)
	}

	return cfg
}

func showHelp() {
	fmt.Println("\nUSAGE:")
	fmt.Println("  fk8cli [options] kubernetes-version")
	fmt.Println("\nEXAMPLE:")
	fmt.Println("  $fk8cli -c 1.34")
	fmt.Println("\nOPTIONS:")
	fmt.Println("  -c  Configure as a control plane node")
	fmt.Println("  -w  Configure as a worker node")
	fmt.Println("  -s  Configure as a single node (control plane + worker)")
	fmt.Println("  -q  Enable quiet output")
	fmt.Println("  -h  Show this help message")
	fmt.Println("\nAt least one of -c, -w, or -s must be specified")
	fmt.Println("The -y flag is required to install Kubernetes and configure the machine as a node")
	fmt.Println("fk8cli must be run with root access (sudo)")
}

// show to-be configuration
func showConfiguration(cfg *Config) {
	fmt.Println("CONFIGURATION:")
	fmt.Println("   Kubernetes version: ", cfg.Tag())
	fmt.Println("   CRI-Tools version:  ", cfg.Tag())
	fmt.Println("   Container Runtime Interface (CRI)")
	fmt.Println("      CRI-O version:   ", cfg.Tag())
	fmt.Println("   Container Runtime")
	fmt.Println("      crun\n")
	fmt.Println("PACKAGES:")
	script.Exec("dnf list " + cfg.Rpms() + " --available").
		Last(len(cfg.rpms)).
		FilterLine(func(line string) string {
			return "   " + line
		}).
		Stdout()
	fmt.Println("\nROLES:")
	if cfg.isControl {
		fmt.Println("   Control node")
	}
	if cfg.isWorker {
		fmt.Println("   Worker node")
	}
}

// show help and exit with exit code
func showHelpAndExit(msg string, exitcode int) {
	prefix := "\nStatus: "
	if exitcode > 0 {
		prefix = "\nError: "
	}
	fmt.Println(prefix, msg, "\n")
	showHelp()
	os.Exit(exitcode)
}

// create array of rpm names to install
func buildRPMList(cfg *Config) {
	k8 := "kubernetes" + cfg.version
	cfg.rpms = append(cfg.rpms, k8, k8+"-client", k8+"-kubeadm", "cri-o"+cfg.version, "crun")
}
