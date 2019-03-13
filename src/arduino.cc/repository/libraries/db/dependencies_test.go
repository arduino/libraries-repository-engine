package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDependencyExtract(t *testing.T) {
	check := func(depDefinition string, name []string, ver []string) {
		dep := extractDependenciesList(depDefinition)
		require.NotNil(t, dep)
		require.Len(t, dep, len(name))
		for i := range name {
			require.Equal(t, name[i], dep[i].Name)
			require.Equal(t, ver[i], dep[i].Version)
		}
	}
	invalid := func(dep string) {
		require.Nil(t, extractDependenciesList(dep))
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
	invalid("MyLib (>=1.2.3)(),AnotherLib, YetAnotherLib (=1.0.0)")
}
