package db

import "encoding/json"

// Version is the type for library versions.
type Version struct {
	version string
}

// Less returns whether the receiver version is lower than the argument.
func (version *Version) Less(other Version) (bool, error) {
	// TODO: apply semantic versioning
	return version.version < other.version, nil
}

// String returns the version in string form.
func (version *Version) String() string {
	return version.version
}

// UnmarshalJSON parses the JSON-encoded argument and stores the result in the receiver.
func (version *Version) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &version.version)
}

// MarshalJSON returns the JSON encoding of the receiver.
func (version *Version) MarshalJSON() ([]byte, error) {
	// Encode version as a string
	return json.Marshal(version.version)
}

// VersionFromString parses a string to a Version object.
func VersionFromString(str string) Version {
	return Version{version: str}
}
