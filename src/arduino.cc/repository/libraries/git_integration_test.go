package libraries

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUpdateLibraryJson(t *testing.T) {
	repos, err := ListRepos("./testdata/git_repos_orgs.txt")

	assert.NoError(t, err)
	assert.NotNil(t, repos)

	for _, repo := range repos {
		repo, err := CloneOrFetch(repo, "/tmp")

		assert.NoError(t, err)
		assert.NotNil(t, repo)

		err = CheckoutLastTag(repo)
		assert.NoError(t, err)
	}

}
