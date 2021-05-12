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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"arduino.cc/repository/libraries/hash"
)

// GithubDownloadRelease downloads GitHub's archive of the release.
func GithubDownloadRelease(repoURL, version string) (string, int64, string, error) {
	tempfile, err := ioutil.TempFile("", "github")
	if err != nil {
		return "", -1, "", err
	}
	defer os.Remove(tempfile.Name())

	zipFileURL := strings.Replace(repoURL, ".git", "", 1) + "/archive/" + version + ".zip"

	err = saveURLIn(zipFileURL, tempfile)
	if err != nil {
		return "", -1, "", err
	}

	info, err := os.Stat(tempfile.Name())
	if err != nil {
		return "", -1, "", err
	}
	size := info.Size()

	checksum, err := hash.Checksum(tempfile.Name())
	if err != nil {
		return "", -1, "", err
	}

	return zipFileURL, size, checksum, nil
}

func saveURLIn(url string, tempfile *os.File) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer tempfile.Close()
	defer resp.Body.Close()

	_, err = io.Copy(tempfile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
