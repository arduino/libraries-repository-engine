package db

import "encoding/json"

type Version struct {
	version string
}

func (x *Version) Less(y Version) (bool, error) {
	// TODO: apply semantic versioning
	return x.version < y.version, nil
}

func (x *Version) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &x.version)
}

func (x *Version) MarshalJSON() ([]byte, error) {
	// Encode version as a string
	return json.Marshal(x.version)
}

func VersionFromString(x string) Version {
	return Version{version: x}
}

// vi:ts=2
