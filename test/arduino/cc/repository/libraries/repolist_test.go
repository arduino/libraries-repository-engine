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
	sort.Sort(libraries.ReposByUrl(repos))

	require.Equal(t, repos[0].Url, "https://github.com/PaulStoffregen/AltSoftSerial.git")
	require.Equal(t, repos[1].Url, "https://github.com/PaulStoffregen/OctoWS2811.git")

	require.Error(t, err)

	error := err.(libraries.GitURLsError)
	require.Equal(t, error.Repos[0].Url, "https://github.com/arduino-libraries")
	require.Equal(t, error.Repos[1].Url, "git@github.com:PaulStoffregen/Audio.git")

}
