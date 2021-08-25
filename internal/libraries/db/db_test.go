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

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testerDB() *DB {
	tDB := DB{
		Libraries: []*Library{
			{
				Name:         "FooLib",
				Repository:   "https://github.com/Bar/FooLib.git",
				SupportLevel: "",
			},
			{
				Name:         "BazLib",
				Repository:   "https://github.com/Bar/BazLib.git",
				SupportLevel: "",
			},
			{
				Name:         "QuxLib",
				Repository:   "https://github.com/Zeb/QuxLib.git",
				SupportLevel: "",
			},
		},
		Releases: []*Release{
			{
				LibraryName:     "FooLib",
				Version:         Version{"1.0.0"},
				Author:          "Barthor",
				Maintainer:      "Bartainer",
				License:         "MIT",
				Sentence:        "asdf",
				Paragraph:       "zxcv",
				Website:         "https://example.com",
				Category:        "Other",
				Architectures:   []string{"avr"},
				Types:           []string{"Contributed"},
				URL:             "http://www.example.com/libraries/github.com/Bar/FooLib-1.0.0.zip",
				ArchiveFileName: "FooLib-1.0.0.zip",
				Size:            123,
				Checksum:        "SHA-256:887f897cfb1818a53652aef39c2a4b8de3c69c805520b2953a562a787b422420",
				Includes:        []string{"FooLib.h"},
				Dependencies: []*Dependency{
					{
						Name:    "BazLib",
						Version: "2.0.0",
					},
				},
				Log: "Some log messages",
			},
			{
				LibraryName:     "BazLib",
				Version:         Version{"2.0.0"},
				Author:          "Barthor",
				Maintainer:      "Bartainer",
				License:         "MIT",
				Sentence:        "asdf",
				Paragraph:       "zxcv",
				Website:         "https://example.com",
				Category:        "Other",
				Architectures:   []string{"avr"},
				Types:           []string{"Contributed"},
				URL:             "http://www.example.com/libraries/github.com/Bar/BazLib-2.0.0.zip",
				ArchiveFileName: "BazLib-2.0.0.zip",
				Size:            123,
				Checksum:        "SHA-256:887f897cfb1818a53652aef39c2a4b8de3c69c805520b2953a562a787b422420",
				Includes:        []string{"BazLib.h"},
				Dependencies:    []*Dependency{},
				Log:             "Some log messages",
			},
			{
				LibraryName:     "BazLib",
				Version:         Version{"2.1.0"},
				Author:          "Barthor",
				Maintainer:      "Bartainer",
				License:         "MIT",
				Sentence:        "asdf",
				Paragraph:       "zxcv",
				Website:         "https://example.com",
				Category:        "Other",
				Architectures:   []string{"avr"},
				Types:           []string{"Contributed"},
				URL:             "http://www.example.com/libraries/github.com/Bar/BazLib-2.1.0.zip",
				ArchiveFileName: "BazLib-2.1.0.zip",
				Size:            123,
				Checksum:        "SHA-256:887f897cfb1818a53652aef39c2a4b8de3c69c805520b2953a562a787b422420",
				Includes:        []string{"BazLib.h"},
				Dependencies:    []*Dependency{},
				Log:             "Some log messages",
			},
			{
				LibraryName:     "FooLib",
				Version:         Version{"1.1.0"},
				Author:          "Barthor",
				Maintainer:      "Bartainer",
				License:         "MIT",
				Sentence:        "asdf",
				Paragraph:       "zxcv",
				Website:         "https://example.com",
				Category:        "Other",
				Architectures:   []string{"avr"},
				Types:           []string{"Contributed"},
				URL:             "http://www.example.com/libraries/github.com/Bar/FooLib-1.1.0.zip",
				ArchiveFileName: "FooLib-1.1.0.zip",
				Size:            123,
				Checksum:        "SHA-256:887f897cfb1818a53652aef39c2a4b8de3c69c805520b2953a562a787b422420",
				Includes:        []string{"FooLib.h"},
				Dependencies: []*Dependency{
					{
						Name:    "BazLib",
						Version: "",
					},
				},
				Log: "Some log messages",
			},
		},
		libraryFile: "some-file.json",
	}

	return &tDB
}

func TestRemoveLibrary(t *testing.T) {
	testDB := testerDB()
	assert.True(t, testDB.HasLibrary("FooLib"))
	assert.True(t, testDB.HasReleaseByNameVersion("FooLib", "1.0.0"))
	err := testDB.RemoveLibrary("FooLib")
	require.NoError(t, err)
	assert.False(t, testDB.HasLibrary("FooLib"))
	assert.False(t, testDB.HasReleaseByNameVersion("FooLib", "1.0.0"))
	assert.False(t, testDB.HasReleaseByNameVersion("FooLib", "1.1.0"))

	assert.True(t, testDB.HasLibrary("QuxLib"))
	err = testDB.RemoveLibrary("QuxLib")
	require.NoError(t, err)
	assert.False(t, testDB.HasLibrary("QuxLib"))

	err = testDB.RemoveLibrary("nonexistent")
	assert.Error(t, err)
}

func TestRemoveReleaseByNameVersion(t *testing.T) {
	testDB := testerDB()
	assert.True(t, testDB.HasReleaseByNameVersion("FooLib", "1.0.0"))
	assert.True(t, testDB.HasReleaseByNameVersion("FooLib", "1.1.0"))
	err := testDB.RemoveReleaseByNameVersion("FooLib", "1.0.0")
	require.NoError(t, err)
	assert.False(t, testDB.HasReleaseByNameVersion("FooLib", "1.0.0"))
	assert.True(t, testDB.HasReleaseByNameVersion("FooLib", "1.1.0"))

	err = testDB.RemoveReleaseByNameVersion("nonexistent", "1.0.0")
	assert.Error(t, err)
	err = testDB.RemoveReleaseByNameVersion("FooLib", "99.99.99")
	assert.Error(t, err)
}

func TestRemoveReleases(t *testing.T) {
	testDB := testerDB()
	assert.True(t, testDB.HasReleaseByNameVersion("FooLib", "1.0.0"))
	err := testDB.RemoveReleases("FooLib")
	require.NoError(t, err)
	assert.False(t, testDB.HasReleaseByNameVersion("FooLib", "1.0.0"))
	assert.False(t, testDB.HasReleaseByNameVersion("FooLib", "1.1.0"))

	err = testDB.RemoveReleases("nonexistent")
	assert.Error(t, err)
}
