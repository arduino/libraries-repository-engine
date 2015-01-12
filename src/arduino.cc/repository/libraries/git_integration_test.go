package libraries

import (
	"testing"
	"github.com/stretchr/testify/require"
	"arduino.cc/repository/libraries/db"
)

func TestUpdateLibraryJson(t *testing.T) {
	repos, err := ListRepos("./testdata/git_only_our_org.txt")

	require.NoError(t, err)
	require.NotNil(t, repos)

	libraryDb := db.Init("./testdata/test_db.json")

	for _, repo := range repos {
		repo, err := CloneOrFetch(repo, "/tmp")

		require.NoError(t, err)
		require.NotNil(t, repo)

		err = CheckoutLastTag(repo)
		require.NoError(t, err)

		library, err := GenerateLibraryFromRepo(repo)
		require.NoError(t, err)
		require.NotNil(t, library)

		err = UpdateLibrary(library, libraryDb)
		require.NoError(t, err)
	}

}
