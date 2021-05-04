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
