package libraries

import (
	"testing"
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/assert"
)

func TestUpdateLibraryJson(t *testing.T) {
	repos, err := libraries.ListRepos("./git_repos_orgs.txt")

	assert.NoError(t, err)
	assert.NotNil(t, repos)

	for _, repo := range repos {
		repo, err := libraries.CloneOrFetch(repo, "/tmp")

		assert.NoError(t, err)
		assert.NotNil(t, repo)

		err = libraries.CheckoutLastTag(repo)
		assert.NoError(t, err)
	}

}
