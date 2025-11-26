/*
Copyright 2025 Bradley G Smith >bradley.g.smith@gmail.com>
SPDX-License-Identifier: MIT

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Functions that handle creating and writing config files
// Assumptions:
// Creates file with contents unique to kubernetes stack and
// not created externally
// Flow:
// check if target exists (assume correct content initially)
// if so exit
// otherwise create temp file locally using go routines
// move to target location via sudo exec
package node

import (
	"errors"
	"io/fs"
	"os"

	"buckaroogeek.com/fk8cli/internal/initialize"
	"buckaroogeek.com/fk8cli/internal/logger"
)

// write to file - if file exists assume identical content
func idempotentWrite(path string, data []byte, cfg *initialize.Config) error {
	// 1. Read existing content
	_, err := os.Open(path)
	if err != nil {
		// If file doesn't exist, create temp file, move to path
		if errors.Is(err, fs.ErrNotExist) {
			var createErr error
			f, createErr := os.CreateTemp("", "fk8cli")
			logger.CheckFatal(createErr)

			defer os.Remove(f.Name())

			// write to temp
			_, createErr = f.Write(data)
			logger.CheckFatal(createErr)

			//copy to path - assumes path directory exists
			cmd := "cp " + f.Name() + " " + path
			createErr = sudoexec(cmd, cfg)
			if createErr != nil {
				return createErr
			}

			// set permissions
			cmd = "chmod 0644 " + path
			return sudoexec(cmd, cfg)
		}

		// Handle other read errors
		logger.CheckFatal(err)
	}

	// file exists
	return nil
}
