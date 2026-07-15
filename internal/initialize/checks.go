/*
Copyright 2025 Bradley G Smith >bradley.g.smith@gmail.com>
SPDX-License-Identifier: MIT

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Sanity check functions
package initialize

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/bitfield/script"
)

// sanity check for Fedora
// exits program with message if not Fedora
func CheckFedora() {
	_, err := script.File("/etc/redhat-release").Match("Fedora").CountLines()
	if err != nil {
		fmt.Println("Fedora is required")
		os.Exit(1)
	}
}

// sanity check for root permissions
func checkRoot() {
	if os.Geteuid() != 0 {
		showHelpAndExit("Root access is required", 1)
	}
}

// sanity check for sudo permissions
func checkSudo() {
	_, err := script.Exec("sudo -n true").String()
	if err != nil {
		showHelpAndExit("Sudo is required", 1)
	}
}

// check swap
// if swap support enabled, add zram-generator to rpm list
// if swap support disabled, fail if zram-generator is present
func checkSwap(cfg *Config) {
	// check for zram-generator, currently installed by default
	pipe := script.Exec("rpm -q zram-generator")
	pipe.Wait()
	zram := pipe.ExitStatus() == 0

	// if zram and cfg.swap identical exit
	if zram == cfg.swap {
		return
	}

	// if zram but swap is false exit
	if zram {
		msg := "Execute 'dnf remove zram-generator' and reboot "
		msg += strconv.FormatBool(zram)
		msg += " "
		msg += strconv.FormatBool(cfg.swap)
		showHelpAndExit(msg, 1)
	}

	cfg.rpms = append(cfg.rpms, "zram-generator")
}

// check for extra arguments
// display args, help, exit
func checkArgs() {
	var errStr string
	if flag.NArg() < 1 {
		errStr = "Target version missing"
	} else if flag.NArg() > 1 {
		errStr = "Too many arguments"
		fmt.Printf("Number of unknown args : %d\n", flag.NArg()-1)
		for i := 1; i < flag.NArg(); i++ {
			fmt.Printf(" %q", flag.Arg(i))
		}
	} else {
		return
	}
	showHelpAndExit(errStr, 1)
}

// check first argument as version - remove 'v' if present
func checkVersion(cfg *Config) {
	cfg.version = strings.TrimLeft(flag.Arg(0), "vV")
}
