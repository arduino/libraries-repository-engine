// This file is part of libraries-repository-engine.
//
// Copyright 2025 ARDUINO SA (http://www.arduino.cc/)
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

package checkregistry

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryValidation(t *testing.T) {
	type testcase struct {
		Name           string
		TestFile       string
		ExpectedResult string
	}
	tests := []testcase{
		{"EmptyArg", "", "registry data file argument testdata is a folder, not a file"},
		{"NonExistentFile", "nonexistent.txt", "while loading registry data file: stat testdata/nonexistent.txt: no such file or directory"},
		//{"InvalidDataFormat", "invalid-data-format.txt", "while loading registry data file: invalid line format (3 fields are required): https://github.com/arduino-libraries/SD.git|Partner;SD"},
		{"InvalidUrlFormat", "invalid-url-format.txt", "while filtering registry data file: Following URL are unknown or unsupported git repos:\nhttps://github.com/arduino-libraries/SD\n"},
		{"MissingType", "no-type.txt", "invalid type '' used by library 'SD'"},
		{"InvalidType", "invalid-type.txt", "invalid type 'foo' used by library 'SD'"},
		{"DuplicateRepoURL", "duplicate-url.txt", "registry data file contains duplicate URLs"},
		{"DuplicateLibName", "duplicate-name.txt", "registry data file contains duplicates of name 'SD'"},
		{"ValidList", "valid.txt", ""},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := runcheck(filepath.Join("testdata", test.TestFile))
			if test.ExpectedResult == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.ExpectedResult)
			}
		})
	}
}
