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

package archive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arduino/libraries-repository-engine/internal/configuration"
	"github.com/arduino/libraries-repository-engine/internal/libraries"
	"github.com/arduino/libraries-repository-engine/internal/libraries/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDataPath string

func init() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testDataPath = filepath.Join(workingDirectory, "testdata")
}

func TestNew(t *testing.T) {
	repository := libraries.Repository{
		FolderPath: "/qux/repos/some-repo",
		URL:        "https://github.com/Foo/Bar.git",
	}
	libraryMetadata := metadata.LibraryMetadata{
		Name:    "Foo Bar",
		Version: "1.2.3",
	}
	config := configuration.Config{
		LibrariesFolder: "/baz/libs/",
		BaseDownloadURL: "https://example/com/libraries/",
	}

	archiveObject, err := New(&repository, &libraryMetadata, &config)
	require.NoError(t, err)
	assert.Equal(t, "/qux/repos/some-repo", archiveObject.SourcePath)
	assert.Equal(t, "Foo_Bar-1.2.3", archiveObject.RootName)
	assert.Equal(t, "Foo_Bar-1.2.3.zip", archiveObject.FileName)
	assert.Equal(t, filepath.Join("/baz/libs/github.com/Foo/Foo_Bar-1.2.3.zip"), archiveObject.Path)
	assert.Equal(t, "https://example/com/libraries/github.com/Foo/Foo_Bar-1.2.3.zip", archiveObject.URL)
}

func TestCreate(t *testing.T) {
	archiveDir := filepath.Join(os.TempDir(), "TestCreateArchiveDir")
	defer os.RemoveAll(archiveDir)
	archivePath := filepath.Join(archiveDir, "TestCreateArchive.zip")

	archiveObject := Archive{
		Path:       archivePath,
		SourcePath: filepath.Join(testDataPath, "gitclones", "SomeRepository"),
		RootName:   "SomeLibrary",
	}

	err := archiveObject.Create()
	require.NoError(t, err, "This test must be run as administrator on Windows to have symlink creation privilege.")

	assert.FileExists(t, archivePath)
	assert.Greater(t, archiveObject.Size, int64(0))
	assert.NotEmpty(t, archiveObject.Checksum)
}
