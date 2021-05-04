package metadata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersionFix(t *testing.T) {
	version, err := VersionToSemverCompliant("1")
	require.NoError(t, err)
	require.Equal(t, "1.0.0", version)

	version, err = VersionToSemverCompliant("1.2")
	require.NoError(t, err)
	require.Equal(t, "1.2.0", version)

	/*
		version, err = VersionToSemverCompliant("1.2b")
		require.NoError(t, err)
		require.Equal(t, "1.2.1", version)

		version, err = VersionToSemverCompliant("05")
		require.NoError(t, err)
		require.Equal(t, "5.0.0", version)
	*/
}
