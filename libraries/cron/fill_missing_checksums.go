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

package cron

import (
	"io"
	"net/http"
	"os"

	"arduino.cc/repository/libraries/hash"
)

/*
FillMissingChecksumsForDownloadArchives checks for missing size and checksum field and fills them
by downloading a copy of the file.
*/
func FillMissingChecksumsForDownloadArchives(URL string, filename string) (int64, string, error) {
	size, err := download(URL, filename)
	if err != nil {
		return 0, "", err
	}

	hash, err := hash.Checksum(filename)
	if err != nil {
		os.Remove(filename)
		return 0, "", err
	}

	return size, hash, nil
}

func download(URL string, filename string) (int64, error) {
	out, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	resp, err := http.Get(URL)
	if err != nil {
		defer os.Remove(out.Name())
		return 0, err
	}
	defer resp.Body.Close()

	size, err := io.Copy(out, resp.Body)
	if err != nil {
		defer os.Remove(out.Name())
		return 0, err
	}

	return size, nil
}

// vi:ts=2
