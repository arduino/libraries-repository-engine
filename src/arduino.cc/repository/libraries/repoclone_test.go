package libraries

import (
	"testing"
	"github.com/stretchr/testify/require"
	"os"
)

func TestCloneRepos(t *testing.T) {
	repo, err := CloneOrFetch("https://github.com/arlibs/Servo.git", "/tmp")

	require.NoError(t, err)
	require.NotNil(t, repo)

	_, err = os.Stat("/tmp/Servo")
	require.NoError(t, err)
}

func TestLastTag(t *testing.T) {
	repo, err := CloneOrFetch("https://github.com/arlibs/Servo.git", "/tmp")

	require.NoError(t, err)
	require.NotNil(t, repo)

	err = CheckoutLastTag(repo)

	require.NoError(t, err)
}
