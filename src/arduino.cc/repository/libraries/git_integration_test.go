package libraries

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestUpdateLibraryJson(t *testing.T) {
	repos, err := ListRepos("./testdata/git_only_our_org.txt")

	require.NoError(t, err)
	require.NotNil(t, repos)

	for _, repo := range repos {
		repo, err := CloneOrFetch(repo, "/tmp")

		require.NoError(t, err)
		require.NotNil(t, repo)

		err = CheckoutLastTag(repo)
		require.NoError(t, err)

		err = UpdateLibrary(repo)
		assert.NoError(t, err)
	}

}
