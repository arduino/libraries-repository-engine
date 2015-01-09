package libraries

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestCloneRepos(t *testing.T) {
	repo, err := CloneOrFetch("https://github.com/arlibs/Servo.git", "/tmp")

	assert.NoError(t, err)
	assert.NotNil(t, repo)

	_, err = os.Stat("/tmp/Servo")
	assert.NoError(t, err)
}

func TestLastTag(t *testing.T) {
	repo, err := CloneOrFetch("https://github.com/arlibs/Servo.git", "/tmp")

	assert.NoError(t, err)
	assert.NotNil(t, repo)

	err = CheckoutLastTag(repo)

	assert.NoError(t, err)
}
