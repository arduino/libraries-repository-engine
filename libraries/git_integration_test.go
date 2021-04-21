package libraries

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"arduino.cc/repository/libraries/db"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/require"
)

func TestUpdateLibraryJson(t *testing.T) {
	repos, err := ListRepos("./testdata/git_only_servo.txt")

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
		err = repoTree.Checkout(&git.CheckoutOptions{Hash: tag.Hash()})
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