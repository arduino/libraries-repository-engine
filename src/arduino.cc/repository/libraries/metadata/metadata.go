/*
A package to handle library.properties metadata.

The functions in this package helps on parsing/validation of
library.properties metadata. All metadata are parsed into a
LibraryMetadata structure.

The source of may be any of the following:
- a github.PullRequest
- a github.RepositoryContent
- a byte[]
*/
package metadata

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/google/go-github/github"
	ini "github.com/vaughan0/go-ini"
)

// Metadata for a library.properties file
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
}

const CATEGORY_UNCATEGORIZED string = "Uncategorized"

func IsValidCategory(category string) bool {
	validCategories := []string{
		"Display",
		"Communication",
		"Signal Input/Output",
		"Sensors",
		"Device Control",
		"Timing",
		"Data Storage",
		"Data Processing",
		"Other",
		CATEGORY_UNCATEGORIZED,
	}
	for _, c := range validCategories {
		if category == c {
			return true
		}
	}
	return false
}

func (library *LibraryMetadata) Validate() []error {
	var errorsAccumulator []error

	// Check author and mainteiner existence
	if library.Author == "" || library.Maintainer == "" {
		errorsAccumulator = append(errorsAccumulator, errors.New("'author' and 'maintainer' fields must be defined"))
	}

	// Check sentence and paragraph and url existence
	if library.Sentence == "" || library.URL == "" {
		errorsAccumulator = append(errorsAccumulator, errors.New("'sentence' and 'url' fields must be defined"))
	}

	newVersion, err := VersionToSemverCompliant(library.Version)
	if err != nil {
		errorsAccumulator = append(errorsAccumulator, err)
	}
	library.Version = newVersion

	// Check if the category is valid and set to "Uncategorized" if not
	if !IsValidCategory(library.Category) {
		library.Category = CATEGORY_UNCATEGORIZED
	}

	return errorsAccumulator
}

// Make a LibraryMetadata by reading library.properties from a github.PullRequest
func ParsePullRequest(gh *github.Client, pull *github.PullRequest) (*LibraryMetadata, error) {
	head := *pull.Head
	headRepo := *head.Repo

	// Get library.properties from pull request HEAD
	getContentOpts := &github.RepositoryContentGetOptions{
		Ref: *head.SHA,
	}
	libPropContent, _, _, err := gh.Repositories.GetContents(*headRepo.Owner.Login, *headRepo.Name, "library.properties", getContentOpts)
	if err != nil {
		return nil, err
	}
	if libPropContent == nil {
		return nil, errors.New("library.properties file not found")
	}
	return ParseRepositoryContent(libPropContent)
}

// Make a LibraryMetadata by reading library.properties from a github.RepositoryContent
func ParseRepositoryContent(content *github.RepositoryContent) (*LibraryMetadata, error) {
	libPropertiesData, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return nil, err
	}
	return Parse(libPropertiesData)
}

// Make a LibraryMetadata by parsing a library.properties file contained in a byte array
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
		} else {
			return ""
		}
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
	}
	return library, nil
}
