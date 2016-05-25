package db

import (
	"arduino.cc/repository/libraries/metadata"
	"strings"
)

func FromLibraryToRelease(library *metadata.LibraryMetadata) *Release {
	architectures := strings.Split(library.Architectures, ",")
	for i, v := range architectures {
		architectures[i] = strings.TrimSpace(v)
	}

	version, err := ParseVersion(library.Version)
	if err != nil {
		panic(err)
	}

	dbRelease := Release{
		LibraryName:   library.Name,
		Version:       version,
		Author:        library.Author,
		Maintainer:    library.Maintainer,
		License:       library.License,
		Sentence:      library.Sentence,
		Paragraph:     library.Paragraph,
		Website:       library.URL, // TODO: Rename "url" field to "website" in library.properties
		Category:      library.Category,
		Architectures: architectures,
		Types:         library.Types,
	}

	return &dbRelease
}
