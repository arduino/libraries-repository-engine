package db

import (
	"strings"

	"arduino.cc/repository/libraries/metadata"
)

func FromLibraryToRelease(library *metadata.LibraryMetadata) *Release {
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
