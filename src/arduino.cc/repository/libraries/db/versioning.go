package db

import "strings"
import "github.com/blang/semver"

func ParseVersion(version string) (semver.Version, error) {
	//version = replaceEndingLetter(version, 'b')

	//version = removeLeading(version, '0')

	versionParts := strings.Split(version, ".")
	for i := len(versionParts); i < 3; i++ {
		versionParts = append(versionParts, "0")
	}
	version = strings.Join(versionParts, ".")

	return semver.Parse(version)
}
