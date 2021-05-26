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
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDataPath string

func init() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testDataPath = filepath.Join(workingDirectory, "testdata")
}

func TestRunArduinoLint(t *testing.T) {
	testTables := []struct {
		testName       string
		folder         string
		official       bool
		reportRegexp   string
		errorAssertion assert.ErrorAssertionFunc
	}{
		{
			"update",
			"Arduino_MKRRGB",
			true,
			"^$",
			assert.NoError,
		},
		{
			"official",
			"Arduino_TestOff",
			true,
			"^$",
			assert.NoError,
		},
		{
			"unofficial",
			"Arduino_Test3rd",
			false,
			"LP012",
			assert.NoError,
		},
		{
			"error",
			"Arduino_TestErr",
			true,
			"LS006",
			assert.Error,
		},
		{
			"warning",
			"Arduino_TestWarn",
			true,
			"LP015",
			assert.NoError,
		},
		{
			"pass",
			"Arduino_TestPass",
			true,
			"^$",
			assert.NoError,
		},
	}

	for _, testTable := range testTables {
		var metadata Repo
		if testTable.official {
			metadata.Types = []string{"Arduino"}
		} else {
			metadata.Types = []string{"Contributed"}
		}
		report, err := RunArduinoLint("", filepath.Join(testDataPath, "libraries", testTable.folder), &metadata)
		assert.Regexp(t, regexp.MustCompile(testTable.reportRegexp), string(report), testTable.testName)
		testTable.errorAssertion(t, err, testTable.testName)
	}
}
