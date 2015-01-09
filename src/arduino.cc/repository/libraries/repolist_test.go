package libraries

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListRepos(t *testing.T) {
	repos, err := ListRepos("./testdata/git_repos_orgs.txt")

	require.NoError(t, err)
	require.Equal(t, len(repos), 3)
}
