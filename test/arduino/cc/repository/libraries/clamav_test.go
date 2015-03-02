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
	require.NoError(t, err)
	defer os.RemoveAll(libraryRepo)

	err = libraries.RunAntiVirus(libraryRepo)
	require.NoError(t, err)
}
