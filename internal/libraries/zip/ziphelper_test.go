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
	"archive/zip"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZip(t *testing.T) {
	zipFile, err := ioutil.TempFile("", "ziphelper*.zip")
	require.NoError(t, err)
	require.NotNil(t, zipFile)
	zipFileName := zipFile.Name()
	require.NoError(t, zipFile.Close())
	require.NoError(t, os.Remove(zipFileName))
	defer os.RemoveAll(zipFileName)

	err = Directory("./testzip", "a_zip", zipFileName)
	require.NoError(t, err)

	zipFileReader, err := zip.OpenReader(zipFileName)
	require.NoError(t, err)

	defer zipFileReader.Close()

	require.Equal(t, 4, len(zipFileReader.File))

	containsName := func(name string) bool {
		for _, file := range zipFileReader.File {
			if file.Name == name {
				return true
			}
		}

		return false
	}
	require.True(t, containsName("a_zip/"))
	require.True(t, containsName("a_zip/testfile.txt"))
	require.True(t, containsName("a_zip/testfolder/"))
	require.True(t, containsName("a_zip/testfolder/testfileinfolder.txt"))
}
