package libraries

import (
	"arduino.cc/repository/libraries"
	"arduino.cc/repository/libraries/db"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"os"
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

	for _, repoURL := range repos {
		repo, err := libraries.CloneOrFetch(repoURL, "/tmp")

		require.NoError(t, err)
		require.NotNil(t, repo)

		defer os.RemoveAll("/tmp/Servo")

		err = libraries.CheckoutLastTag(repo)
		require.NoError(t, err)

		library, err := libraries.GenerateLibraryFromRepo(repo.Workdir())
		require.NoError(t, err)
		require.NotNil(t, library)

		release := db.FromLibraryToRelease(library, "http://www.example.com/")

		err = libraries.UpdateLibrary(release, libraryDb)
		require.NoError(t, err)

		err = libraries.ZipRepo(repo.Workdir(), library, librariesRepo)
		require.NoError(t, err)
	}
}
