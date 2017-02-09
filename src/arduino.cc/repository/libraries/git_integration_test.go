package libraries

import (
	"io/ioutil"
	"os"
	"testing"

	"arduino.cc/repository/libraries/db"
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
		r, err := CloneOrFetch(repo.Url, "/tmp")

		require.NoError(t, err)
		require.NotNil(t, r)

		defer os.RemoveAll(r.FolderPath)

		err = CheckoutLastTag(r)
		require.NoError(t, err)

		library, err := GenerateLibraryFromRepo(r)
		require.NoError(t, err)
		require.NotNil(t, library)

		zipFolderName := ZipFolderName(library)

		release := db.FromLibraryToRelease(library)

		zipFilePath, err := ZipRepo(r.FolderPath, librariesRepo, zipFolderName)
		require.NoError(t, err)
		require.NotEmpty(t, zipFilePath)

		err = UpdateLibrary(release, libraryDb)
		require.NoError(t, err)

	}
}
