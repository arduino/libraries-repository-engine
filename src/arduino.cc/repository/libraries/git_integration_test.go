package libraries

import (
	"testing"
	"github.com/stretchr/testify/require"
	"arduino.cc/repository/libraries/db"
	"io/ioutil"
)

func TestUpdateLibraryJson(t *testing.T) {
	repos, err := ListRepos("./testdata/git_only_servo.txt")

	require.NoError(t, err)
	require.NotNil(t, repos)

	librariesRepo, err := ioutil.TempDir("", "libraries")
	require.NoError(t, err)

	libraryDb := db.Init("./testdata/test_db.json")

	for _, repoURL := range repos {
		repo, err := CloneOrFetch(repoURL, "/tmp")

		require.NoError(t, err)
		require.NotNil(t, repo)

		err = CheckoutLastTag(repo)
		require.NoError(t, err)

		library, err := GenerateLibraryFromRepo(repo.Workdir())
		require.NoError(t, err)
		require.NotNil(t, library)

		err = UpdateLibrary(library, libraryDb)
		require.NoError(t, err)

		err = ZipRepo(repo.Workdir(), library, librariesRepo)
		require.NoError(t, err)
	}
}
