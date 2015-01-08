package libraries

import (
	"testing"
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestCloneRepos(t *testing.T) {
	repo, err := libraries.CloneOrFetch("https://github.com/arlibs/Servo.git", "/tmp")

	assert.NoError(t, err)
	assert.NotNil(t, repo)

	_, err = os.Stat("/tmp/Servo")
	assert.NoError(t, err)
}

func TestLastTag(t *testing.T) {
	repo, err := libraries.CloneOrFetch("https://github.com/arlibs/Servo.git", "/tmp")

	assert.NoError(t, err)
	assert.NotNil(t, repo)

	err = libraries.CheckoutLastTag(repo)

	assert.NoError(t, err)
}
