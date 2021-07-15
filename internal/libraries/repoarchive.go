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

	"github.com/arduino/libraries-repository-engine/internal/libraries/metadata"
	"github.com/arduino/libraries-repository-engine/internal/libraries/zip"
)

// ZipRepo creates a ZIP archive of the repo folder and returns its path.
func ZipRepo(repoFolder string, baseFolder string, zipFolderName string) (string, error) {
	err := os.MkdirAll(baseFolder, os.FileMode(0755))
	if err != nil {
		return "", err
	}
	absoluteFileName := filepath.Join(baseFolder, zipFolderName+".zip")
	if err := zip.Directory(repoFolder, zipFolderName, absoluteFileName); err != nil {
		os.Remove(absoluteFileName)
		return "", err
	}

	return absoluteFileName, nil
}

// ZipFolderName returns the name to use for the folder.
func ZipFolderName(library *metadata.LibraryMetadata) string {
	pattern := regexp.MustCompile("[^a-zA-Z0-9]")
	return pattern.ReplaceAllString(library.Name, "_") + "-" + library.Version
}
