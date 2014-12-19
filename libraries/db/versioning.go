package db

import "errors"
import "encoding/json"

type Version struct {
	version string
}

func (x *Version) Less(y *Version) (bool, error) {
	if x == nil || y == nil {
		return false, errors.New("Invalid version in comparison")
	}
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

func VersionFromString(x *string) *Version {
	if x == nil {
		return nil
	}
	var v Version
	v.version = *x
	return &v
}

// vi:ts=2
