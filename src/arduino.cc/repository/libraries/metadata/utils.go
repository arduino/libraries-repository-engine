package metadata

import (
	"strings"

	"github.com/blang/semver"
)

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
