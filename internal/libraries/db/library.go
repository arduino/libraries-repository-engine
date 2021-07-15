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

package db

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arduino/libraries-repository-engine/internal/libraries/metadata"
)

// FromLibraryToRelease extract a Release from LibraryMetadata. LibraryMetadata must be
// validated before running this function.
func FromLibraryToRelease(library *metadata.LibraryMetadata) *Release {
	deps, _ := ExtractDependenciesList(library.Depends)
	dbRelease := Release{
		LibraryName:   library.Name,
		Version:       VersionFromString(library.Version),
		Author:        library.Author,
		Maintainer:    library.Maintainer,
		License:       library.License,
		Sentence:      library.Sentence,
		Paragraph:     library.Paragraph,
		Website:       library.URL, // TODO: Rename "url" field to "website" in library.properties
		Category:      library.Category,
		Architectures: extractStringList(library.Architectures),
		Types:         library.Types,
		Includes:      extractStringList(library.Includes),
		Dependencies:  deps,
	}

	return &dbRelease
}

func extractStringList(value string) []string {
	split := strings.Split(value, ",")
	res := []string{}
	for _, s := range split {
		s := strings.TrimSpace(s)
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}

var re = regexp.MustCompile("^([a-zA-Z0-9](?:[a-zA-Z0-9._\\- ]*[a-zA-Z0-9])?) *(?: \\(([^()]*)\\))?$")

// ExtractDependenciesList extracts dependencies from the "depends" field of library.properties
func ExtractDependenciesList(depends string) ([]*Dependency, error) {
	deps := []*Dependency{}
	depends = strings.TrimSpace(depends)
	if depends == "" {
		return deps, nil
	}
	for _, dep := range strings.Split(depends, ",") {
		dep = strings.TrimSpace(dep)
		if dep == "" {
			return nil, fmt.Errorf("invalid dep: %s", dep)
		}
		matches := re.FindAllStringSubmatch(dep, -1)
		if matches == nil {
			return nil, fmt.Errorf("invalid dep: %s", dep)
		}
		deps = append(deps, &Dependency{
			Name:    matches[0][1],
			Version: matches[0][2],
		})
	}
	return deps, nil
}
