package db

// Output structure used to generate library_index.json file
type indexOutput struct {
	Libraries []indexLibrary `json:"libraries"`
}

// Output structure used to generate library_index.json file
type indexLibrary struct {
	LibraryName      string   `json:"name"`
	Version          Version  `json:"version"`
	Author           string   `json:"author"`
	Maintainer       string   `json:"maintainer"`
	License          string   `json:"license,omitempty"`
	Sentence         string   `json:"sentence"`
	Paragraph        string   `json:"paragraph,omitempty"`
	Website          string   `json:"website,omitempty"`
	Category         string   `json:"category,omitempty"`
	Architectures    []string `json:"architectures"`
	Types            []string `json:"types,omitempty"`
	ProvidesIncludes []string `json:"provides_includes,omitempty"`

	URL             string `json:"url"`
	ArchiveFileName string `json:"archiveFileName"`
	Size            int64  `json:"size"`
	Checksum        string `json:"checksum"`

	SupportLevel string `json:"supportLevel,omitempty"`
}

// Generate an object that once JSON-marshaled produces a json
// file suitable for the library installer (i.e. produce a valid
// library_index.json file)
func (db *DB) OutputLibraryIndex() (interface{}, error) {
	libraries := make([]indexLibrary, 0, len(db.Libraries))

	for _, lib := range db.Libraries {
		libraryReleases := db.FindReleasesOfLibrary(lib)

		for _, libraryRelease := range libraryReleases {
			// Skip malformed release
			if libraryRelease.Size == 0 || libraryRelease.Checksum == "" {
				continue
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
				ProvidesIncludes: libraryRelease.Includes,
			})
		}

	}

	index := indexOutput{
		Libraries: libraries,
	}
	return &index, nil
}
