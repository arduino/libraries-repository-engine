package metadata

import (
	"github.com/blang/semver"
	"strings"
)

func VersionToSemverCompliant(version string) (string, error) {
	//version = replaceEndingLetter(version, 'b')

	//version = removeLeading(version, '0')

	versionParts := strings.Split(version, ".")
	if len(versionParts) < 3 {
		for i := len(versionParts); i < 3; i++ {
			versionParts = insert(versionParts, i, "0")
		}
	}
	version = strings.Join(versionParts, ".")

	newVersion, err := semver.Parse(version)
	if err != nil {
		return "", err
	}
	return newVersion.String(), nil
}

func replaceEndingLetter(version string, letter uint8) string {
	if version[len(version)-1] == letter {
		version = version[:len(version)-1]+".1"
	}
	return version
}

func removeLeading(version string, letter uint8) string {
	for version[0] == letter {
		version = version[1:]
	}
	return version
}

func insert(original []string, position int, value string) []string {
	//we'll grow by 1
	target := make([]string, len(original)+1)

	//copy everything up to the position
	copy(target, original[:position])

	//set the new value at the desired position
	target[position] = value

	//copy everything left over
	copy(target[position+1:], original[position:])

	return target
}
