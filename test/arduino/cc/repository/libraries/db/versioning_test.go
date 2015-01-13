package db

import (
	"arduino.cc/repository/libraries/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLess(t *testing.T) {
	v1 := db.VersionFromString("1.0")
	v2 := db.VersionFromString("2.0")

	res, err := v1.Less(v2)
	require.NoError(t, err)
	require.True(t, res)
}

func TestUnmarshalJSON(t *testing.T) {
	v1 := db.Version{}
	err := v1.UnmarshalJSON([]byte("\"1.0\""))
	require.NoError(t, err)
	require.Equal(t, "1.0", v1.String())
}

func TestMarshalJSON(t *testing.T) {
	v1 := db.VersionFromString("1.0")
	bytes, err := v1.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "\"1.0\"", string(bytes))
}
