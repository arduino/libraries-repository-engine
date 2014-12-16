package db

// Output structure used to generate library_index.json file
type indexOutput struct {
	Libraries []indexLibrary `json:"libraries"`
}

type indexLibrary struct {
	LibraryName   *string  `json:"name"`
	Version       *Version `json:"version"`
	Author        *string  `json:"author"`
	Maintainer    *string  `json:"maintainer"`
	License       *string  `json:"license,omitempty"`
	Sentence      *string  `json:"sentence"`
	Paragraph     *string  `json:"paragraph,omitempty"`
	Website       *string  `json:"website,omitempty"`
	Category      *string  `json:"category,omitempty"`
	Architectures []string `json:"architectures,omitempty"`

	URL      *string `json:"url"`
	Size     uint64  `json:"size"`
	Checksum *string `json:"checksum"`
}

func (db *DB) OutputLibraryIndex() (interface{}, error) {
	libraries := make([]indexLibrary, 0, len(db.Libraries))

	for _, lib := range db.Libraries {
		latest, err := db.FindLatestReleaseOfLibrary(&lib)
		if err != nil {
			return nil, err
		}

		// Skip libraries without release
		if latest == nil {
			continue
		}

		// Skip malformed release
		if latest.Size == 0 || latest.Checksum == nil {
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
			URL:           latest.URL,
			Size:          latest.Size,
			Checksum:      latest.Checksum,
		})
	}

	index := &indexOutput{
		Libraries: libraries,
	}
	return index, nil
}

// vi:ts=2
