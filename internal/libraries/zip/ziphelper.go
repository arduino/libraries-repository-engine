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

package zip

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/arduino/libraries-repository-engine/internal/libraries/file"
)

// Directory creates a new zip archive that contains a copy of "rootFolder" into "zipFile".
// Inside the archive "rootFolder" will be renamed to "zipRootFolderName".
func Directory(rootFolder string, zipRootFolderName string, zipFile string) error {
	checks := func(path string, info os.FileInfo, err error) error {
		info, err = os.Lstat(path)
		if err != nil {
			return err
		}
		if (info.Mode() & os.ModeSymlink) != 0 {
			dest, _ := os.Readlink(path)
			return fmt.Errorf("Symlink not allowed: %s -> %s", path, dest)
		}
		if file.IsSCCS(info.Name()) {
			return filepath.SkipDir
		}
		return nil
	}

	rootFolder, err := filepath.Abs(rootFolder)
	if err != nil {
		return err
	}

	if err := filepath.Walk(rootFolder, checks); err != nil {
		return err
	}

	tmpdir, err := ioutil.TempDir("", "ziphelper")
	if err != nil {
		return fmt.Errorf("creating temp dir for zip archive: %s", err)
	}
	defer os.RemoveAll(tmpdir)
	if err := os.Symlink(rootFolder, filepath.Join(tmpdir, zipRootFolderName)); err != nil {
		return fmt.Errorf("creating temp dir for zip archive: %s", err)
	}

	args := []string{"-r", zipFile, zipRootFolderName, "-x", ".*", "-x", "*/.*"}
	for sccs := range file.SCCSFiles {
		args = append(args, "-x")
		args = append(args, "*/"+sccs+"/*")
		args = append(args, "-x")
		args = append(args, sccs+"/*")
	}

	zipCmd := exec.Command("zip", args...)
	zipCmd.Dir = tmpdir
	if _, err := zipCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("archiving into zip: %s", err)
	}
	return nil
}
