package libraries

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListRepos(t *testing.T) {
	repos, err := ListRepos("./testdata/git_repos_orgs.txt")

	assert.NoError(t, err)
	assert.Equal(t, len(repos), 3)
}
