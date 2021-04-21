package metadata

import (
	"strings"

	"github.com/blang/semver"
)

// VersionToSemverCompliant transforms a truncated version to a semver compliant version,
// for example "1.0" is converted to "1.0.0".
func VersionToSemverCompliant(version string) (string, error) {
	versionParts := len(strings.Split(version, "."))
	for versionParts < 3 {
		versionParts++
		version += ".0"
	}

	newVersion, err := semver.Parse(version)
	if err != nil {
		return "", err
	}
	return newVersion.String(), nil
}
