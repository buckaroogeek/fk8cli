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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/bitfield/script"
)

/* Sanity checks, parse flags and arguements */
func ParseFlags() *Config {
	cfg := NewConfig()

	flag.BoolVarP(&cfg.isControl, "control", "c", cfg.isControl, "Configure as a control plane node")
	flag.BoolVarP(&cfg.isWorker, "worker", "w", cfg.isWorker, "Configure as a worker node")
	flag.BoolVarP(&cfg.isSingle, "single", "s", cfg.isSingle, "Configure as a single node (control plane + worker)")
	flag.BoolVarP(&cfg.localrpms, "local", "l", cfg.localrpms, "Install rpms from the local directory")
	flag.BoolVarP(&cfg.isQuiet, "quiet", "q", cfg.isQuiet, "Enable verbose output")
	flag.BoolVarP(&cfg.swap, "swap", "a", cfg.swap, "Enable swap support")
	flag.BoolVarP(&cfg.taint, "taint", "t", cfg.taint, "Set taint on control plane node")
	flag.BoolVarP(&cfg.isGo, "yes", "y", cfg.isGo, "Proceed with installation")

	flag.Usage = showHelp
	flag.Parse()

	// check for no args and flags
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		showHelpAndExit("No options or version were provided", 0)
	}

	// Create logfile name
	buildLogFileName(cfg)

	// check swap status
	checkSwap(cfg)

	// check root access
	checkSudo()

	// user must not combine -s with -c or -w
	if cfg.isSingle {
		if cfg.isControl || cfg.isWorker {
			showHelpAndExit("Cannot combine -s with -c and/or -w", 1)
		}
		// toggle control and worker to true given isSingle is set
		cfg.isControl = true
		cfg.isWorker = true
		cfg.taint = true
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

	// Retrieve user name
	getUserName(cfg)

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
	fmt.Println("  $fk8cli -c -y 1.36 - installs a v1.36 control plane without taint")
	fmt.Println("\nOPTIONS:")
	fmt.Println("  -y  --yes      Execute the installation")
	fmt.Println("\n  -c  --control  Configure as a control plane node")
	fmt.Println("  -w  --worker   Configure as a worker node")
	fmt.Println("  -s  --single   Configure as a single node (control plane + worker)")
	fmt.Println("  -l  --locale   Install rpms from local directory")
	fmt.Println("                 Local rpms installed instead of rpms from repo")
	fmt.Println("  -q  --quiet    Enable quiet output")
	fmt.Println("\n  -a  --swap     Add support for swap")
	fmt.Println("  -t  --taint    Set taint on control plane node")
	fmt.Println("                 Taint set automatically on single node")
	fmt.Println("\n  -h  --help     Show this help message")
	fmt.Println("\nNotes:")
	fmt.Println("* At least one of -c, -w, or -s must be specified")
	fmt.Println("* The -y flag is required to proceed with installation and configuration")
	fmt.Println("* The fk8cli user must have sudo")
	fmt.Println("* Run dnf update and reboot before using this utility")
}

// show to-be configuration
func showConfiguration(cfg *Config) {
	fmt.Println("USER:", cfg.User())
	fmt.Println("CONFIGURATION:")
	fmt.Println("   Kubernetes version: ", cfg.Tag())
	fmt.Println("   CRI-Tools version:  ", cfg.Tag())
	fmt.Println("   Container Runtime Interface (CRI)")
	fmt.Println("      CRI-O version:   ", cfg.Tag())
	fmt.Println("   Container Runtime")
	fmt.Println("      crun\n")

	// dnf output
	fmt.Println("REPO PACKAGES:")
	// fmt.Println("   dnf string: " + cfg.Rpms())
	script.Exec("sudo dnf list " + cfg.Rpms()).FilterLine(func(line string) string {
		return "   " + line
	}).Stdout()

	// local rpms if any
	if cfg.LocalRpms() {
		fmt.Println("\nLOCAL RPMS:")
		for _, rpmname := range cfg.rpmfiles {
			fmt.Println("   " + rpmname)
		}
	}
	fmt.Println("\nROLES:")
	if cfg.isControl {
		var withstr string
		withstr = " without"
		if cfg.GetTaint() {
			withstr = " with"
		}
		fmt.Println("   Control plane" + withstr + " taint")
	}
	if cfg.isWorker {
		fmt.Println("   Worker")
	}
	fmt.Println("\nLOG: ", cfg.FileName())
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

// build log file name
func buildLogFileName(cfg *Config) {
	t := time.Now()
	cfg.filename = "fk8cli_" + t.Format(time.DateOnly) + ".log"
}

// create array of rpm names to install
func buildRPMList(cfg *Config) {

	// build list of local rpm files if flagged
	// remove duplicates from reporpms list
	if cfg.localrpms {
		list, err := filepath.Glob("*.rpm")
		if err != nil {
			log.Fatal(err)
		}
		if len(list) > 0 {
			// replace rpmfiles map using list size
			cfg.rpmfiles = make(map[string]string, len(list))

			// build map of package names extracted from rpm name
			for _, spec := range list {
				name, err := script.Exec("rpm -qp --qf '%{NAME}' " + spec).String()
				if err != nil {
					log.Fatal(err)
				}
				// fmt.Println(name, " ", spec)
				cfg.rpmfiles[name] = spec
			}
		} else {
			err := errors.New("Local rpms selected for install but none found")
			log.Fatal(err)
		}
		// fmt.Println(cfg.rpmfiles)
	}

	// build list of rpms to install from repo, filtering by
	// list of rpms to be installed from filesystem
	// kubernetes rpms
	k8 := "kubernetes" + cfg.version
	cfg.AddRPM(k8)
	cfg.AddRPM(k8 + "-client")
	cfg.AddRPM(k8 + "-kubeadm")

	// CRI and crictl
	cfg.AddRPM("cri-o" + cfg.version)
	cfg.AddRPM("cri-tools" + cfg.version)

	// container runtime
	cfg.AddRPM("crun")
}

// retrieve user name
func getUserName(cfg *Config) {
	user, err := script.Exec("logname").String()
	if err != nil {
		log.Fatal(err)
	}
	cfg.user = strings.TrimSpace(user)
}
