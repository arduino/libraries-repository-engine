package db

// Output structure used to generate library_index.json file
type indexOutput struct {
	Libraries []indexLibrary `json:"libraries"`
}

// Output structure used to generate library_index.json file
type indexLibrary struct {
	LibraryName   string   `json:"name"`
	Version       Version  `json:"version"`
	Author        string   `json:"author"`
	Maintainer    string   `json:"maintainer"`
	License       string   `json:"license,omitempty"`
	Sentence      string   `json:"sentence"`
	Paragraph     string   `json:"paragraph,omitempty"`
	Website       string   `json:"website,omitempty"`
	Category      string   `json:"category,omitempty"`
	Architectures []string `json:"architectures,omitempty"`

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
		latest, err := db.FindLatestReleaseOfLibrary(lib)
		if err != nil {
			return nil, err
		}

		// Skip libraries without release
		if latest == nil {
			continue
		}

		// Skip malformed release
		if latest.Size == 0 || latest.Checksum == "" {
			continue
		}

		// Copy db.Library into db.indexLibrary
		libraries = append(libraries, indexLibrary{
			LibraryName:   latest.LibraryName,
			Version:       latest.Version,
			Author:        latest.Author,
			Maintainer:    latest.Maintainer,
			License:       latest.License,
			Sentence:      latest.Sentence,
			Paragraph:     latest.Paragraph,
			Website:       latest.Website,
			Category:      latest.Category,
			Architectures: latest.Architectures,

			ArchiveFileName: latest.ArchiveFileName,
			URL:             latest.URL,
			Size:            latest.Size,
			Checksum:        latest.Checksum,

			SupportLevel: lib.SupportLevel,
		})
	}

	index := indexOutput{
		Libraries: libraries,
	}
	return &index, nil
}
