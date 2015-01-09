package libraries

import (
	"testing"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"arduino.cc/repository/libraries/metadata"
	"fmt"
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

		bytes, err := ioutil.ReadFile(repo.Workdir() + "library.properties")
		require.NoError(t, err)

		library, err := metadata.Parse(bytes)
		require.NoError(t, err)

		errs := library.Validate()
		require.Empty(t, errs)

		fmt.Println(library)
	}

}
