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
	isControl bool
	isGo      bool
	isQuiet   bool
	isSingle  bool
	isWorker  bool
	rpms      []string
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
		version:   "none",
	}

	return cfg
}

// Methods

// Return list of rpms as single string
func (cfg *Config) Rpms() string {
	return strings.Join(cfg.rpms, " ")
}

// Return target version for k8s
func (cfg *Config) Tag() string {
	return "v" + cfg.version
}

// Return target version for k8s
func (cfg *Config) Version() string {
	return cfg.version
}
