/*
Copyright 2025 Bradley G Smith >bradley.g.smith@gmail.com>
SPDX-License-Identifier: MIT

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Initiates cli parser, then calls Run to launch
// RunAction

package fk8cli

import (
	"buckaroogeek.com/fk8cli/internal/initialize"
	// "buckaroogeek.com/fk8cli/internal/logger"
	"buckaroogeek.com/fk8cli/internal/node"
)

func Main() {
	//sanity check for Fedora
	initialize.CheckFedora()

	//parse flags
	cfg := initialize.ParseFlags()

	node.Install(cfg)
}
