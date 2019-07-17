/*
Package metadata handles library.properties metadata.

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
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"arduino.cc/repository/libraries/db"

	"github.com/google/go-github/github"
	ini "github.com/vaughan0/go-ini"
)

// LibraryMetadata contains metadata for a library.properties file
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
	Includes      string
	Depends       string
}

const categoryUcategorized string = "Uncategorized"

var validCategories = []string{
	"Display",
	"Communication",
	"Signal Input/Output",
	"Sensors",
	"Device Control",
	"Timing",
	"Data Storage",
	"Data Processing",
	"Other",
	categoryUcategorized,
}

// IsValidCategory checks if category is a valid category
func IsValidCategory(category string) bool {
	for _, c := range validCategories {
		if category == c {
			return true
		}
	}
	return false
}

// Validate checks LibraryMetadata for errors, returns an array of the errors found
func (library *LibraryMetadata) Validate() []error {
	var errorsAccumulator []error

	// Check lib name
	if !IsValidLibraryName(library.Name) {
		errorsAccumulator = append(errorsAccumulator, errors.New("Invalid 'name' field: "+library.Name))
	}

	// Check author and maintainer existence
	if library.Author == "" {
		errorsAccumulator = append(errorsAccumulator, errors.New("'author' field must be defined"))
	}
	if library.Maintainer == "" {
		library.Maintainer = library.Author
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
		library.Category = categoryUcategorized
	}

	// Check if 'depends' field is correctly written
	if !IsValidDependency(library.Depends) {
		errorsAccumulator = append(errorsAccumulator, errors.New("Invalid 'depends' field: "+library.Depends))
	}
	return errorsAccumulator
}

// IsValidLibraryName checks if a string is a valid library name
func IsValidLibraryName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if name[0] == '-' || name[0] == '_' || name[0] == ' ' {
		return false
	}
	for _, char := range name {
		if !strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-. ", string(char)) {
			return false
		}
	}
	return true
}

// IsValidDependency checks if the `depends` field of library.properties is correctly formatted
func IsValidDependency(depends string) bool {
	_, err := db.ExtractDependenciesList(depends)
	return err == nil
}

// ParsePullRequest makes a LibraryMetadata by reading library.properties from a github.PullRequest
func ParsePullRequest(gh *github.Client, pull *github.PullRequest) (*LibraryMetadata, error) {
	head := *pull.Head
	headRepo := *head.Repo

	// Get library.properties from pull request HEAD
	getContentOpts := &github.RepositoryContentGetOptions{
		Ref: *head.SHA,
	}
	libPropContent, _, _, err := gh.Repositories.GetContents(context.TODO(), *headRepo.Owner.Login, *headRepo.Name, "library.properties", getContentOpts)
	if err != nil {
		return nil, err
	}
	if libPropContent == nil {
		return nil, errors.New("library.properties file not found")
	}
	return ParseRepositoryContent(libPropContent)
}

// ParseRepositoryContent makes a LibraryMetadata by reading library.properties from a github.RepositoryContent
func ParseRepositoryContent(content *github.RepositoryContent) (*LibraryMetadata, error) {
	libPropertiesData, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return nil, err
	}
	return Parse(libPropertiesData)
}

// Parse makes a LibraryMetadata by parsing a library.properties file contained in a byte array
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
		}
		return ""
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
		Includes:      get("includes"),
		Depends:       get("depends"),
	}
	return library, nil
}
