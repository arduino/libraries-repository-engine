package db

import "encoding/json"

type Version struct {
	version string
}

func (version *Version) Less(other Version) (bool, error) {
	// TODO: apply semantic versioning
	return version.version < other.version, nil
}

func (version *Version) String() string {
	return version.version
}

func (version *Version) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &version.version)
}

func (version *Version) MarshalJSON() ([]byte, error) {
	// Encode version as a string
	return json.Marshal(version.version)
}

func VersionFromString(str string) Version {
	return Version{version: str}
}
