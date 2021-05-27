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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"arduino.cc/repository/internal/libraries/db"
	"arduino.cc/repository/internal/libraries/gitutils"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/require"
)

func TestUpdateLibraryJson(t *testing.T) {
	repos, err := ListRepos("./testdata/git_test_repo.txt")

	require.NoError(t, err)
	require.NotNil(t, repos)

	librariesRepo, err := ioutil.TempDir("", "libraries")
	require.NoError(t, err)
	defer os.RemoveAll(librariesRepo)

	libraryDb := db.Init("./testdata/test_db.json")
	defer os.RemoveAll("./testdata/test_db.json")

	for _, repo := range repos {
		subfolder, err := repo.AsFolder()
		require.NoError(t, err)

		r, err := CloneOrFetch(repo, filepath.Join("/tmp", subfolder))
		require.NoError(t, err)
		require.NotNil(t, r)

		defer os.RemoveAll(r.FolderPath)

		tags, err := r.Repository.Tags()
		require.NoError(t, err)
		tag, err := tags.Next()
		require.NoError(t, err)

		repoTree, err := r.Repository.Worktree()
		require.NoError(t, err)
		// Annotated tags have their own hash, different from the commit hash, so the tag must be resolved before checkout
		resolvedTag, err := gitutils.ResolveTag(tag, r.Repository)
		require.NoError(t, err)
		err = repoTree.Checkout(&git.CheckoutOptions{Hash: *resolvedTag, Force: true})
		require.NoError(t, err)
		err = repoTree.Clean(&git.CleanOptions{Dir: true})
		require.NoError(t, err)

		library, err := GenerateLibraryFromRepo(r)
		require.NoError(t, err)
		require.NotNil(t, library)

		zipFolderName := ZipFolderName(library)

		release := db.FromLibraryToRelease(library)

		zipFilePath, err := ZipRepo(r.FolderPath, librariesRepo, zipFolderName)
		require.NoError(t, err)
		require.NotEmpty(t, zipFilePath)

		err = UpdateLibrary(release, r.URL, libraryDb)
		require.NoError(t, err)

	}
}
