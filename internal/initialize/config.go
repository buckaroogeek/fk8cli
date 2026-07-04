/*
Copyright 2025 Bradley G Smith >bradley.g.smith@gmail.com>
SPDX-License-Identifier: MIT

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Configuration struct populated by convention and by
// command line flags (see flags)
package initialize

import (
	"strings"
)

// configuration for node set up
type Config struct {
	filename  string
	isControl bool
	isGo      bool
	isQuiet   bool
	isSingle  bool
	isWorker  bool
	localrpms bool
	rpmfiles  []string
	rpms      []string
	taint     bool
	user      string
	version   string
}

// Config constructor with explicit defaults
func NewConfig() *Config {
	cfg := &Config{
		isControl: false,
		isGo:      false,
		isWorker:  false,
		isSingle:  false,
		isQuiet:   false,
		localrpms: false,
		taint:     false,
		version:   "none",
	}

	return cfg
}

// Methods

// Return log file name
func (cfg *Config) FileName() string {
	return cfg.filename
}

// Return verbose boolean
func (cfg *Config) IsVerbose() bool {
	return !cfg.isQuiet
}

// Return localrpms boolean
func (cfg *Config) LocalRpms() bool {
	return cfg.localrpms
}

// Return list of local rpmfiles as single string
func (cfg *Config) Rpmfiles() string {
	return strings.Join(cfg.rpmfiles, " ")
}

// Return list of rpms as single string
func (cfg *Config) Rpms() string {
	return strings.Join(cfg.rpms, " ")
}

// Return taint boolean
func (cfg *Config) SetTaint() bool {
	return cfg.taint
}

// Return target version for k8s
func (cfg *Config) Tag() string {
	return "v" + cfg.version
}

// Return user name
func (cfg *Config) User() string {
	return cfg.user
}

// Return target version for k8s
func (cfg *Config) Version() string {
	return cfg.version
}
