// This file is part of libraries-repository-engine.
//
// Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

/*
Package metadata handles library.properties metadata.

The functions in this package helps on parsing/validation of
library.properties metadata. All metadata are parsed into a
LibraryMetadata structure.
*/
package metadata

import (
	"bytes"

	"github.com/arduino/arduino-cli/arduino/libraries"
	ini "github.com/vaughan0/go-ini"
	semver "go.bug.st/relaxed-semver"
)

// LibraryMetadata contains metadata for a library.properties file
type LibraryMetadata struct {
	Name          string
	Version       string
	Author        string
	Maintainer    string
	License       string
	Sentence      string
	Paragraph     string
	URL           string
	Architectures string
	Category      string
	Types         []string
	Includes      string
	Depends       string
}

// Parse makes a LibraryMetadata by parsing a library.properties file contained in a byte array
func Parse(propertiesData []byte) (*LibraryMetadata, error) {
	// Create an io.Reader from []bytes
	reader := bytes.NewReader(propertiesData)
	// Use go-ini to decode contents
	properties, err := ini.Load(reader)
	if err != nil {
		return nil, err
	}
	get := func(key string) string {
		value, ok := properties.Get("", key)
		if ok {
			return value
		}
		return ""
	}
	library := &LibraryMetadata{
		Name:          get("name"),
		Version:       get("version"),
		Author:        get("author"),
		Maintainer:    get("maintainer"),
		Sentence:      get("sentence"),
		Paragraph:     get("paragraph"),
		License:       get("license"),
		URL:           get("url"),
		Architectures: get("architectures"),
		Category:      get("category"),
		Includes:      get("includes"),
		Depends:       get("depends"),
	}

	library.normalize()

	return library, nil
}

// normalize normalizes library metadata.
func (library *LibraryMetadata) normalize() {
	library.Version = normalizeVersion(library.Version)
	library.Category = normalizeCategory(library.Category)
}

// normalizeVersion converts "relaxed semver" to semver-compliant versions.
func normalizeVersion(version string) string {
	versionObject, err := semver.Parse(version)
	if err != nil {
		// Enforcement is handled by Arduino Lint.
		return version
	}

	versionObject.Normalize()
	return versionObject.String()
}

// normalizeCategory restricts category values to the allowed list.
func normalizeCategory(category string) string {
	if !libraries.ValidCategories[category] {
		return "Uncategorized"
	}

	return category
}
