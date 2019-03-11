package metadata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataChecks(t *testing.T) {
	require.True(t, IsValidLibraryName("asjdh qweqwe - qweqwe_213123"))
	require.False(t, IsValidLibraryName("-asdasda"))
	require.False(t, IsValidLibraryName("_"))
	require.False(t, IsValidLibraryName("  asdasd"))
	require.False(t, IsValidLibraryName("asd$"))

	require.True(t, IsValidDependency("asdas"))
	require.True(t, IsValidDependency("asdas asdadasd,123213"))
	require.False(t, IsValidDependency("asdas asdadasd,,123213"))
	require.True(t, IsValidDependency(""))
	require.False(t, IsValidDependency("_123123,asdasd"))
	require.False(t, IsValidDependency("435regf,asdkwjqwe,_ger,213123"))
}
