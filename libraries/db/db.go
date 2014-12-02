package db

import "bytes"
import "github.com/vaughan0/go-ini"

// Metadata for a library.properties file
type LibraryMetadata struct {
	Name          *string
	Version       *string
	Author        *string
	Maintainer    *string
	Sentence      *string
	Paragraph     *string
	URL           *string
	Architectures *string
}

// Create a Library by reading library.properties from a byte array
func ParseLibraryProperties(propertiesData []byte) (*LibraryMetadata, error) {
	// Create an io.Reader from []bytes
	reader := bytes.NewReader(propertiesData)
	// Use go-ini to decode contents
	properties, err := ini.Load(reader)
	if err != nil {
		return nil, err
	}
	get := func(key string) *string {
		value, ok := properties.Get("", key)
		if ok {
			return &value
		} else {
			return nil
		}
	}
	library := &LibraryMetadata{
		Name:          get("name"),
		Version:       get("version"),
		Author:        get("author"),
		Maintainer:    get("maintainer"),
		Sentence:      get("sentence"),
		Paragraph:     get("paragraph"),
		URL:           get("url"),
		Architectures: get("architectures"),
	}
	return library, nil
}

// vi:ts=2
