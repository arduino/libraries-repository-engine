package db

import (
	"arduino.cc/repository/libraries/metadata"
	"strings"
)

func FromLibraryToRelease(library *metadata.LibraryMetadata, baseDownloadURL string) *Release {
	architectures := strings.Split(library.Architectures, ",")
	for i, v := range architectures {
		architectures[i] = strings.TrimSpace(v)
	}

	archiveFileName := library.Name + "-" + library.Version + ".zip"
	dbRelease := Release{
		LibraryName:     library.Name,
		Version:         VersionFromString(library.Version),
		Author:          library.Author,
		Maintainer:      library.Maintainer,
		License:         library.License,
		Sentence:        library.Sentence,
		Paragraph:       library.Paragraph,
		Website:         library.URL, // TODO: Rename "url" field to "website" in library.properties
		Category:        library.Category,
		Architectures:   architectures,
		URL:             baseDownloadURL + archiveFileName,
		ArchiveFileName: archiveFileName,
	}

	return &dbRelease
}
