/*
Copyright 2025 Bradley G Smith >bradley.g.smith@gmail.com>
SPDX-License-Identifier: MIT

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Sets up logging - log (eventually) to file and to
// stdout if verbose
//
// logging notes
// routine msgs to standard out track progress, any errors
// log file is a record of what is installed and when, duplicates standard out

package logger

import (
	"fmt"
	"log"
)

// check function to call log.Fatal
func CheckFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Log(token string) {
	fmt.Println("***", token)
}
