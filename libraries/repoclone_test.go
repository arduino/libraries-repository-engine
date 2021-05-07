package libraries

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloneRepos(t *testing.T) {
	meta := &Repo{URL: "https://github.com/arduino-libraries/Servo.git"}

	subfolder, err := meta.AsFolder()
	require.NoError(t, err)

	repo, err := CloneOrFetch(meta, filepath.Join("/tmp", subfolder))

	require.NoError(t, err)
	require.NotNil(t, repo)
	require.Equal(t, "/tmp/github.com/arduino-libraries/Servo", repo.FolderPath)

	defer os.RemoveAll(repo.FolderPath)

	_, err = os.Stat(repo.FolderPath)
	require.NoError(t, err)
}
