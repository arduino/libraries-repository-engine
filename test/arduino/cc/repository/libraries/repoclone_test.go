package libraries

import (
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCloneRepos(t *testing.T) {
	repo, err := libraries.CloneOrFetch("https://github.com/arlibs/Servo.git", "/tmp")

	require.NoError(t, err)
	require.NotNil(t, repo)

	defer os.RemoveAll("/tmp/Servo")

	_, err = os.Stat("/tmp/Servo")
	require.NoError(t, err)
}

func TestLastTag(t *testing.T) {
	repo, err := libraries.CloneOrFetch("https://github.com/arlibs/Servo.git", "/tmp")

	require.NoError(t, err)
	require.NotNil(t, repo)

	defer os.RemoveAll("/tmp/Servo")

	err = libraries.CheckoutLastTag(repo)

	require.NoError(t, err)
}
