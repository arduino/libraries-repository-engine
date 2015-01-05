package libraries

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"arduino.cc/repository/libraries"
)

func TestListRepos(t *testing.T) {
	repos, err := libraries.ListRepos("./git_repos_orgs.txt")

	assert.NoError(t, err)
	assert.Equal(t, len(repos), 3)

}
