package db

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestLess(t *testing.T) {
	v1 := VersionFromString("1.0")
	v2 := VersionFromString("2.0")

	res, err := v1.Less(v2)
	require.NoError(t, err)
	require.True(t, res)
}

func TestUnmarshalJSON(t *testing.T) {
	v1 := Version{}
	err := v1.UnmarshalJSON([]byte("\"1.0\""))
	require.NoError(t, err)
	require.Equal(t, "1.0", v1.version)
}

func TestMarshalJSON(t *testing.T) {
	v1 := VersionFromString("1.0")
	bytes, err := v1.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "\"1.0\"", string(bytes))
}
