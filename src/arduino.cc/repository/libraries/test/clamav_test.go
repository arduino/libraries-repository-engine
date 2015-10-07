package libraries

import (
	"arduino.cc/repository/libraries"
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/require"
	"os"
)

func TestRunClamAV(t *testing.T) {
	libraryRepo, err := ioutil.TempDir("", "library")
	os.Chmod(libraryRepo, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(libraryRepo)

	err = libraries.RunAntiVirus(libraryRepo)
	require.NoError(t, err)
}
