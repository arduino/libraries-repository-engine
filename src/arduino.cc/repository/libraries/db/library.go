package db

import (
	"fmt"
	"regexp"
	"strings"

	"arduino.cc/repository/libraries/metadata"
)

// FromLibraryToRelease extract a Release from LibraryMetadata. LibraryMetadata must be
// validated before running this function.
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
		Dependencies:  extractDependenciesList(library.Depends),
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

func extractDependenciesList(depends string) []*Dependency {
	deps := []*Dependency{}
	for _, dep := range strings.Split(depends, ",") {
		dep = strings.TrimSpace(dep)
		if dep == "" {
			continue
		}
		matches := re.FindAllStringSubmatch(dep, -1)
		if matches == nil {
			fmt.Println("invalid dep:", dep)
			return nil
		}
		deps = append(deps, &Dependency{
			Name:    matches[0][1],
			Version: matches[0][2],
		})
	}
	return deps
}
