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

	"github.com/stretchr/testify/require"
)

func TestDependencyExtract(t *testing.T) {
	check := func(depDefinition string, name []string, ver []string) {
		dep, err := ExtractDependenciesList(depDefinition)
		require.NoError(t, err)
		require.NotNil(t, dep)
		require.Len(t, dep, len(name))
		for i := range name {
			require.Equal(t, name[i], dep[i].Name)
			require.Equal(t, ver[i], dep[i].Version)
		}
	}
	invalid := func(depends string) {
		dep, err := ExtractDependenciesList(depends)
		require.Nil(t, dep)
		require.Error(t, err)
	}
	invalid("-invalidname")
	invalid("_invalidname")
	check("ciao", []string{"ciao"}, []string{""})
	check("MyLib (>=1.2.3)", []string{"MyLib"}, []string{">=1.2.3"})
	check("MyLib (>=1.2.3),AnotherLib, YetAnotherLib (=1.0.0)",
		[]string{"MyLib", "AnotherLib", "YetAnotherLib"},
		[]string{">=1.2.3", "", "=1.0.0"})
	invalid("MyLib (>=1.2.3)()")
	invalid("MyLib (>=1.2.3),_aaaa")
	invalid("MyLib,,AnotherLib")
	invalid("MyLib (>=1.2.3)(),AnotherLib, YetAnotherLib (=1.0.0)")
	check("Arduino Uno WiFi Dev Ed Library, LoRa Node (^2.1.2)",
		[]string{"Arduino Uno WiFi Dev Ed Library", "LoRa Node"},
		[]string{"", "^2.1.2"})
	check("Arduino Uno WiFi Dev Ed Library   ,   LoRa Node    (^2.1.2)",
		[]string{"Arduino Uno WiFi Dev Ed Library", "LoRa Node"},
		[]string{"", "^2.1.2"})
	check("Arduino_OAuth, ArduinoHttpClient (<0.3.0), NonExistentLib",
		[]string{"Arduino_OAuth", "ArduinoHttpClient", "NonExistentLib"},
		[]string{"", "<0.3.0", ""})
	check("", []string{}, []string{})
}
