package libraries

import (
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestRunClamAV(t *testing.T) {
	libraryRepo, err := ioutil.TempDir("", "library")
	os.Chmod(libraryRepo, 0755)
	require.NoError(t, err)
	defer os.RemoveAll(libraryRepo)

	err = libraries.RunAntiVirus(libraryRepo)
	require.NoError(t, err)
}
