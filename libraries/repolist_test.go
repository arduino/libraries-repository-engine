package libraries

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListRepos(t *testing.T) {
	repos, err := ListRepos("./testdata/git_test_repos.txt")

	require.Equal(t, 11, len(repos))

	require.Equal(t, "https://github.com/PaulStoffregen/OctoWS2811.git", repos[0].Url)
	require.Equal(t, "https://github.com/PaulStoffregen/AltSoftSerial.git", repos[1].Url)

	require.Equal(t, "https://github.com/Cheong2K/ble-sdk-arduino.git", repos[2].Url)
	require.Equal(t, "https://github.com/arduino-libraries/Bridge.git", repos[3].Url)
	require.Equal(t, "https://github.com/adafruit/Adafruit_ADS1X15.git", repos[4].Url)
	require.Equal(t, "https://github.com/adafruit/Adafruit_ADXL345.git", repos[5].Url)
	require.Equal(t, "https://github.com/adafruit/Adafruit_AHRS.git", repos[6].Url)
	require.Equal(t, "https://github.com/adafruit/Adafruit_AM2315.git", repos[7].Url)
	require.Equal(t, "https://github.com/arduino-libraries/Scheduler.git", repos[8].Url)
	require.Equal(t, "https://github.com/arduino-libraries/SD.git", repos[9].Url)
	require.Equal(t, "https://github.com/arduino-libraries/Servo.git", repos[10].Url)
	require.Error(t, err)

	error, ok := err.(GitURLsError)
	require.True(t, ok)
	require.Equal(t, "https://github.com/arduino-libraries", error.Repos[0].Url)
	require.Equal(t, "git@github.com:PaulStoffregen/Audio.git", error.Repos[1].Url)
}

func TestRepoFolderPathDetermination(t *testing.T) {
	repo := &Repo{Url: "https://github.com/arduino-libraries/Servo.git"}
	f, err := repo.AsFolder()
	require.NoError(t, err)
	require.Equal(t, "github.com/arduino-libraries/Servo", f)

	repo = &Repo{Url: "https://bitbucket.org/bjoern/arduino_osc"}
	f, err = repo.AsFolder()
	require.NoError(t, err)
	require.Equal(t, "bitbucket.org/bjoern/arduino_osc", f)
}
