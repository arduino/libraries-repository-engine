// This file is part of libraries-repository-engine.
//
// Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package libraries

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var empty struct{}

var officialTypes = map[string]struct{}{
	"Arduino": empty,
}

func official(metadata *Repo) bool {
	for _, libraryType := range metadata.Types {
		_, isOfficial := officialTypes[libraryType]
		if isOfficial {
			return true
		}
	}
	return false
}

// RunArduinoLint runs Arduino Lint on the library and returns the report in the event of error or warnings.
func RunArduinoLint(arduinoLintPath string, folder string, metadata *Repo) ([]byte, error) {
	if arduinoLintPath == "" {
		// Assume Arduino Lint is installed under PATH.
		arduinoLintPath = "arduino-lint"
	}

	JSONReportFolder, err := ioutil.TempDir("", "arduino-lint-report-")
	if err != nil {
		panic(err)
	}
	JSONReportPath := filepath.Join(JSONReportFolder, "report.json")
	defer os.RemoveAll(JSONReportPath)

	// See: https://arduino.github.io/arduino-lint/latest/commands/arduino-lint/
	cmd := exec.Command(
		arduinoLintPath,
		"--compliance=permissive",
		"--format=text",
		"--project-type=library",
		"--recursive=false",
		"--report-file="+JSONReportPath,
		folder,
	)
	// See: https://arduino.github.io/arduino-lint/latest/#environment-variables
	cmd.Env = modifyEnv(os.Environ(), "ARDUINO_LINT_LIBRARY_MANAGER_INDEXING", "true")
	cmd.Env = modifyEnv(cmd.Env, "ARDUINO_LINT_OFFICIAL", fmt.Sprintf("%t", official(metadata)))

	textReport, lintErr := cmd.CombinedOutput()
	if lintErr != nil {
		return textReport, lintErr
	}

	// Read report.
	rawJSONReport, err := ioutil.ReadFile(JSONReportPath)
	if err != nil {
		panic(err)
	}
	var JSONReport map[string]interface{}
	if err := json.Unmarshal(rawJSONReport, &JSONReport); err != nil {
		panic(err)
	}

	// Check warning count.
	reportSummary := JSONReport["summary"].(map[string]interface{})
	warningCount := reportSummary["warningCount"].(float64)

	// Report should be displayed when there are warnings.
	if warningCount > 0 {
		return textReport, lintErr
	}

	// No warnings.
	return nil, nil
}
