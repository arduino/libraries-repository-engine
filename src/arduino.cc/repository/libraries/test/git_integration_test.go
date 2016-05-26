package libraries

import (
	"arduino.cc/repository/libraries"
	"arduino.cc/repository/libraries/db"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestUpdateLibraryJson(t *testing.T) {
	repos, err := libraries.ListRepos("./testdata/git_only_servo.txt")

	require.NoError(t, err)
	require.NotNil(t, repos)

	librariesRepo, err := ioutil.TempDir("", "libraries")
	require.NoError(t, err)
	defer os.RemoveAll(librariesRepo)

	libraryDb := db.Init("./testdata/test_db.json")
	defer os.RemoveAll("./testdata/test_db.json")

	for _, repo := range repos {
		r, err := libraries.CloneOrFetch(repo.Url, "/tmp")

		require.NoError(t, err)
		require.NotNil(t, r)

		defer os.RemoveAll(r.FolderPath)

		err = libraries.CheckoutLastTag(r)
		require.NoError(t, err)

		library, err := libraries.GenerateLibraryFromRepo(r.FolderPath, repo)
		require.NoError(t, err)
		require.NotNil(t, library)

		zipFolderName := libraries.ZipFolderName(library)

		release := db.FromLibraryToRelease(library)

		zipFilePath, err := libraries.ZipRepo(r.FolderPath, librariesRepo, zipFolderName)
		require.NoError(t, err)
		require.NotEmpty(t, zipFilePath)

		err = libraries.UpdateLibrary(release, libraryDb)
		require.NoError(t, err)

	}
}
