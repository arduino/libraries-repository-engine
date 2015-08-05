package libraries

import (
	"testing"
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/assert"
)

func TestBadFileFinderDotDevelopment(t *testing.T) {
	err := libraries.FailIfHasUndesiredFiles("./testdata/lib_with_forbidden_file")
	assert.Error(t, err)
}

func TestBadFileFinderExe(t *testing.T) {
	err := libraries.FailIfHasUndesiredFiles("./testdata/lib_with_exe")
	assert.Error(t, err)
}

func TestBadFileFinderValid(t *testing.T) {
	err := libraries.FailIfHasUndesiredFiles("./testdata/lib_valid")
	assert.NoError(t, err)
}

