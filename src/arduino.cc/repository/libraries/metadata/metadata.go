package metadata

import "bytes"
import "encoding/base64"
import "errors"
import "github.com/vaughan0/go-ini"
import "github.com/cmaglie/go-github/github"

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
}

func (library *LibraryMetadata) Validate() []error {
	var errorsAccumulator []error

	// Check author and mainteiner existence
	if library.Author == nil || library.Maintainer == nil {
		errorsAccumulator = append(errorsAccumulator, errors.New("'author' and 'maintainer' fields must be defined"))
	}

	// Check sentence and paragraph and url existence
	if library.Sentence == nil || library.Paragraph == nil || library.URL == nil {
		errorsAccumulator = append(errorsAccumulator, errors.New("'sentence', 'paragraph' and 'url' fields must be defined"))
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
