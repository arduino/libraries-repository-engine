package libraries

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloneRepos(t *testing.T) {
	repo, err := CloneOrFetch("https://github.com/arduino-libraries/Servo.git", "/tmp")

	require.NoError(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "/tmp/github.com/arduino-libraries/Servo", repo.FolderPath)

	defer os.RemoveAll(repo.FolderPath)

	_, err = os.Stat(repo.FolderPath)
	require.NoError(t, err)
}
