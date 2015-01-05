package libraries

import (
	"testing"
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/assert"
)

func TestCloneRepos(t *testing.T) {
	err := libraries.CloneOrPull("https://github.com/PaulStoffregen/OctoWS2811.git", "/tmp")

	assert.NoError(t, err)
}
