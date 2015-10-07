package metadata

import (
	"testing"
	"github.com/stretchr/testify/require"
	"arduino.cc/repository/libraries/metadata"
)

func TestVersionFix(t *testing.T) {
	var version string
	var err error

	version, err = metadata.VersionToSemverCompliant("1.0")

	require.NoError(t, err)
	require.Equal(t, "1.0.0", version)

	version, err = metadata.VersionToSemverCompliant("1.2")

	require.NoError(t, err)
	require.Equal(t, "1.2.0", version)

	/*
	version, err = metadata.VersionToSemverCompliant("1.2b")

	require.NoError(t, err)
	require.Equal(t, "1.2.1", version)

	version, err = metadata.VersionToSemverCompliant("05")

	require.NoError(t, err)
	require.Equal(t, "5.0.0", version)
	*/
}
