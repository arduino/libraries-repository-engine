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
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

// FailIfHasUndesiredFiles returns an error if the folder contains any undesired files.
func FailIfHasUndesiredFiles(folder string) error {
	err := failIfContainsForbiddenFileInRoot(folder)
	if err != nil {
		return err
	}
	return failIfContainsExes(folder)
}

// ForbiddenFiles is the names of the forbidden files.
var ForbiddenFiles = []string{".development"}

func failIfContainsForbiddenFileInRoot(folder string) error {
	for _, file := range ForbiddenFiles {
		if _, err := os.Stat(filepath.Join(folder, file)); err == nil {
			return errors.New(file + " file found, skipping")
		}
	}

	return nil
}

// Patterns is the file patterns of executables.
var Patterns = []string{"*.exe"}

func failIfContainsExes(folder string) error {
	for _, pattern := range Patterns {
		cmd := exec.Command("find", folder, "-type", "f", "-name", pattern)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		if len(string(output)) > 0 {
			return errors.New(pattern + " files found, skipping")
		}
	}
	return nil
}
