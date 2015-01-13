package libraries

import (
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestListRepos(t *testing.T) {
	repos, err := libraries.ListRepos("./testdata/git_test_repos.txt")

	require.Equal(t, len(repos), 2)
	sort.Strings(repos)

	require.Equal(t, repos[0], "https://github.com/PaulStoffregen/AltSoftSerial.git")
	require.Equal(t, repos[1], "https://github.com/PaulStoffregen/OctoWS2811.git")

	require.Error(t, err)

	error := err.(libraries.GitURLsError)
	require.Equal(t, error.GitURLs[0], "https://github.com/arlibs")
	require.Equal(t, error.GitURLs[1], "git@github.com:PaulStoffregen/Audio.git")

}
