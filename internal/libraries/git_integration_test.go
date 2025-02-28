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

package libraries_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arduino/libraries-repository-engine/internal/configuration"
	"github.com/arduino/libraries-repository-engine/internal/libraries"
	"github.com/arduino/libraries-repository-engine/internal/libraries/archive"
	"github.com/arduino/libraries-repository-engine/internal/libraries/db"
	"github.com/arduino/libraries-repository-engine/internal/libraries/gitutils"
	"github.com/stretchr/testify/require"
)

func TestUpdateLibraryJson(t *testing.T) {
	repos, err := libraries.ListRepos("./testdata/git_test_repo.txt")

	require.NoError(t, err)
	require.NotNil(t, repos)

	librariesRepo, err := os.MkdirTemp("", "libraries")
	require.NoError(t, err)
	defer os.RemoveAll(librariesRepo)

	libraryDb := db.Init("./testdata/test_db.json")
	defer os.RemoveAll("./testdata/test_db.json")

	for _, repo := range repos {
		subfolder, err := repo.AsFolder()
		require.NoError(t, err)

		r, err := libraries.CloneOrFetch(repo, filepath.Join("/tmp", subfolder))
		require.NoError(t, err)
		require.NotNil(t, r)

		defer os.RemoveAll(r.FolderPath)

		tags, err := r.Repository.Tags()
		require.NoError(t, err)
		tag, err := tags.Next()
		require.NoError(t, err)

		err = gitutils.CheckoutTag(r.Repository, tag)
		require.NoError(t, err)

		library, err := libraries.GenerateLibraryFromRepo(r)
		require.NoError(t, err)
		require.NotNil(t, library)

		config := configuration.Config{LibrariesFolder: librariesRepo}
		archiveData, err := archive.New(r, library, &config)
		require.NoError(t, err)

		release := db.FromLibraryToRelease(library)

		err = archiveData.Create()
		require.NoError(t, err)
		require.NotEmpty(t, archiveData.Path)

		err = libraries.UpdateLibrary(release, r.URL, libraryDb)
		require.NoError(t, err)

	}
}
