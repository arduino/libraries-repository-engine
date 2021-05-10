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

// Output structure used to generate library_index.json file
type indexOutput struct {
	Libraries []indexLibrary `json:"libraries"`
}

// Output structure used to generate library_index.json file
type indexLibrary struct {
	LibraryName      string             `json:"name"`
	Version          Version            `json:"version"`
	Author           string             `json:"author"`
	Maintainer       string             `json:"maintainer"`
	License          string             `json:"license,omitempty"`
	Sentence         string             `json:"sentence"`
	Paragraph        string             `json:"paragraph,omitempty"`
	Website          string             `json:"website,omitempty"`
	Category         string             `json:"category,omitempty"`
	Architectures    []string           `json:"architectures"`
	Types            []string           `json:"types,omitempty"`
	Repository       string             `json:"repository,omitempty"`
	ProvidesIncludes []string           `json:"providesIncludes,omitempty"`
	Dependencies     []*indexDependency `json:"dependencies,omitempty"`
	URL              string             `json:"url"`
	ArchiveFileName  string             `json:"archiveFileName"`
	Size             int64              `json:"size"`
	Checksum         string             `json:"checksum"`

	SupportLevel string `json:"supportLevel,omitempty"`
}

type indexDependency struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// OutputLibraryIndex generates an object that once JSON-marshaled produces a json
// file suitable for the library installer (i.e. produce a valid library_index.json file)
func (db *DB) OutputLibraryIndex() (interface{}, error) {
	libraries := make([]indexLibrary, 0, len(db.Libraries))

	for _, lib := range db.Libraries {
		libraryReleases := db.FindReleasesOfLibrary(lib)

		for _, libraryRelease := range libraryReleases {
			// Skip malformed release
			if libraryRelease.Size == 0 || libraryRelease.Checksum == "" {
				continue
			}

			deps := []*indexDependency{}
			for _, dep := range libraryRelease.Dependencies {
				deps = append(deps, &indexDependency{
					Name:    dep.Name,
					Version: dep.Version,
				})
			}

			// Copy db.Library into db.indexLibrary
			libraries = append(libraries, indexLibrary{
				LibraryName:      libraryRelease.LibraryName,
				Version:          libraryRelease.Version,
				Author:           libraryRelease.Author,
				Maintainer:       libraryRelease.Maintainer,
				License:          libraryRelease.License,
				Sentence:         libraryRelease.Sentence,
				Paragraph:        libraryRelease.Paragraph,
				Website:          libraryRelease.Website,
				Category:         lib.LatestCategory,
				Architectures:    libraryRelease.Architectures,
				Types:            libraryRelease.Types,
				ArchiveFileName:  libraryRelease.ArchiveFileName,
				URL:              libraryRelease.URL,
				Size:             libraryRelease.Size,
				Checksum:         libraryRelease.Checksum,
				SupportLevel:     lib.SupportLevel,
				Repository:       lib.Repository,
				ProvidesIncludes: libraryRelease.Includes,
				Dependencies:     deps,
			})
		}

	}

	index := indexOutput{
		Libraries: libraries,
	}
	return &index, nil
}
