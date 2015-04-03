package libraries

import (
	"testing"
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/require"
)

func TestGithubDownload(t *testing.T) {
	url, size, checksum, err := libraries.GithubDownloadRelease("https://github.com/arduino-libraries/Audio.git", "1.0.0")

	require.NoError(t, err)
	require.Equal(t, url, "https://github.com/arduino-libraries/Audio/archive/1.0.0.zip")
	require.Equal(t, checksum, "SHA-256:ccafda9702f98a1332adfdc19cf242479cf715129fdcf47928db81adf2bf2077")
	require.Equal(t, size, 7314)

}
