package libraries

import (
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGithubDownload(t *testing.T) {
	url, size, checksum, err := libraries.GithubDownloadRelease("https://github.com/arduino-libraries/Audio.git", "1.0.0")

	require.NoError(t, err)
	require.Equal(t, url, "https://github.com/arduino-libraries/Audio/archive/1.0.0.zip")
	require.Equal(t, checksum, "SHA-256:a4da301186904b0f95ea691681b40867bd5a1fe608963a79a7c0a2d45f80a320")
	require.Equal(t, size, int64(7314))

}
