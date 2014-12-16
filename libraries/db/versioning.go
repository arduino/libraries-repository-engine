package db

import "errors"

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
	x.version = string(data)
	return nil
}

func (x *Version) MarshalJSON() ([]byte, error) {
	r := make([]byte, len(x.version))
	copy(r[:], x.version)
	return r, nil
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
