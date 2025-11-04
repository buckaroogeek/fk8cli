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
package node

import (
	"fmt"

	"buckaroogeek.com/fk8cli/internal/initialize"
)

// Manage installation process
func Install(cfg *initialize.Config) {
	fmt.Println("K8S: ", cfg.Version())
}
