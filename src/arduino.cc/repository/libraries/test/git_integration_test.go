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
		repoFolder, err := libraries.CloneOrFetch(repo.Url, "/tmp")

		require.NoError(t, err)
		require.NotNil(t, repoFolder)

		defer os.RemoveAll(repoFolder)

		err = libraries.CheckoutLastTag(repoFolder)
		require.NoError(t, err)

		library, err := libraries.GenerateLibraryFromRepo(repoFolder, repo)
		require.NoError(t, err)
		require.NotNil(t, library)

		zipFolderName := libraries.ZipFolderName(library)

		release := db.FromLibraryToRelease(library, "http://www.example.com/", zipFolderName+".zip")

		zipFilePath, err := libraries.ZipRepo(repoFolder, librariesRepo, zipFolderName)
		require.NoError(t, err)
		require.NotEmpty(t, zipFilePath)

		err = libraries.UpdateLibrary(release, libraryDb)
		require.NoError(t, err)

	}
}
