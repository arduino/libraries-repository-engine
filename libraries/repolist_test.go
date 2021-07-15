// This file is part of libraries-repository-engine.
//
// Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package libraries

import (
	"testing"

	"github.com/arduino/libraries-repository-engine/internal/libraries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadRepoListFromFile(t *testing.T) {
	_, err := LoadRepoListFromFile("./testdata/nonexistent.txt")
	assert.Error(t, err, "Attempt to load non-existent registry data file")

	repos, err := LoadRepoListFromFile("./testdata/git_test_repos.txt")
	require.NoError(t, err)

	reposAssertion := []*Repo{
		{
			URL:         "https://github.com/arduino-libraries",
			Types:       []string{"Arduino"},
			LibraryName: "libraries",
		},
		{
			URL:         "git@github.com:PaulStoffregen/Audio.git",
			Types:       []string{"Contributed"},
			LibraryName: "Audio",
		},
		{
			URL:         "https://github.com/PaulStoffregen/OctoWS2811.git",
			Types:       []string{"Arduino", "Contributed"},
			LibraryName: "OctoWS2811",
		},
		{
			URL:         "https://github.com/PaulStoffregen/AltSoftSerial.git",
			Types:       []string{"Contributed"},
			LibraryName: "AltSoftSerial",
		},
		{
			URL:         "https://github.com/Cheong2K/ble-sdk-arduino.git",
			Types:       []string{"Contributed"},
			LibraryName: "ble-sdk-arduino",
		},
		{
			URL:         "https://github.com/arduino-libraries/Bridge.git",
			Types:       []string{"Contributed"},
			LibraryName: "Bridge",
		},
		{
			URL:         "https://github.com/adafruit/Adafruit_ADS1X15.git",
			Types:       []string{"Recommended"},
			LibraryName: "Adafruit_ADS1X15",
		},
		{
			URL:         "https://github.com/adafruit/Adafruit_ADXL345.git",
			Types:       []string{"Recommended"},
			LibraryName: "Adafruit_ADXL345",
		},
		{
			URL:         "https://github.com/adafruit/Adafruit_AHRS.git",
			Types:       []string{"Recommended"},
			LibraryName: "Adafruit_AHRS",
		},
		{
			URL:         "https://github.com/adafruit/Adafruit_AM2315.git",
			Types:       []string{"Recommended"},
			LibraryName: "Adafruit_AM2315",
		},
		{
			URL:         "https://github.com/arduino-libraries/Scheduler.git",
			Types:       []string{"Arduino"},
			LibraryName: "Scheduler",
		},
		{
			URL:         "https://github.com/arduino-libraries/SD.git",
			Types:       []string{"Arduino"},
			LibraryName: "SD",
		},
		{
			URL:         "https://github.com/arduino-libraries/Servo.git",
			Types:       []string{"Arduino"},
			LibraryName: "Servo",
		},
	}

	assert.Equal(t, reposAssertion, repos)
}

func TestListRepos(t *testing.T) {
	repos, err := ListRepos("./testdata/git_test_repos.txt")

	require.Equal(t, 11, len(repos))

	require.Equal(t, "https://github.com/PaulStoffregen/OctoWS2811.git", repos[0].URL)
	require.Equal(t, "https://github.com/PaulStoffregen/AltSoftSerial.git", repos[1].URL)

	require.Equal(t, "https://github.com/Cheong2K/ble-sdk-arduino.git", repos[2].URL)
	require.Equal(t, "https://github.com/arduino-libraries/Bridge.git", repos[3].URL)
	require.Equal(t, "https://github.com/adafruit/Adafruit_ADS1X15.git", repos[4].URL)
	require.Equal(t, "https://github.com/adafruit/Adafruit_ADXL345.git", repos[5].URL)
	require.Equal(t, "https://github.com/adafruit/Adafruit_AHRS.git", repos[6].URL)
	require.Equal(t, "https://github.com/adafruit/Adafruit_AM2315.git", repos[7].URL)
	require.Equal(t, "https://github.com/arduino-libraries/Scheduler.git", repos[8].URL)
	require.Equal(t, "https://github.com/arduino-libraries/SD.git", repos[9].URL)
	require.Equal(t, "https://github.com/arduino-libraries/Servo.git", repos[10].URL)
	require.Error(t, err)

	error, ok := err.(libraries.GitURLsError)
	require.True(t, ok)
	require.Equal(t, "https://github.com/arduino-libraries", error.Repos[0].URL)
	require.Equal(t, "git@github.com:PaulStoffregen/Audio.git", error.Repos[1].URL)
}
