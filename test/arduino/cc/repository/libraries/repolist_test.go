package libraries

import (
	"arduino.cc/repository/libraries"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListRepos(t *testing.T) {
	repos, err := libraries.ListRepos("./testdata/git_test_repos.txt")

	require.Equal(t, len(repos), 11)

	require.Equal(t, repos[0].Url, "https://github.com/PaulStoffregen/OctoWS2811.git")
	require.Equal(t, repos[1].Url, "https://github.com/PaulStoffregen/AltSoftSerial.git")

	require.Equal(t, repos[2].Url, "https://github.com/Cheong2K/ble-sdk-arduino.git")
	require.Equal(t, repos[3].Url, "https://github.com/arduino-libraries/Bridge.git")
	require.Equal(t, repos[4].Url, "https://github.com/adafruit/Adafruit_ADS1X15.git")
	require.Equal(t, repos[5].Url, "https://github.com/adafruit/Adafruit_ADXL345.git")
	require.Equal(t, repos[6].Url, "https://github.com/adafruit/Adafruit_AHRS.git")
	require.Equal(t, repos[7].Url, "https://github.com/adafruit/Adafruit_AM2315.git")
	require.Equal(t, repos[8].Url, "https://github.com/arduino-libraries/Scheduler.git")
	require.Equal(t, repos[9].Url, "https://github.com/arduino-libraries/SD.git")
	require.Equal(t, repos[10].Url, "https://github.com/arduino-libraries/Servo.git")

	require.Error(t, err)

	error := err.(libraries.GitURLsError)
	require.Equal(t, error.Repos[0].Url, "https://github.com/arduino-libraries")
	require.Equal(t, error.Repos[1].Url, "git@github.com:PaulStoffregen/Audio.git")

}
