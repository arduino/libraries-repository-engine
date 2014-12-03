package db

import "bytes"
import "encoding/json"
import "github.com/vaughan0/go-ini"
import "io"
import "os"

// The libraries DB
type DB struct {
	Libraries []Library "libraries"
}

// A library
type Library struct {
	Name          string   "name"
	Version       string   "version"
	Author        string   "author"
	Maintainer    string   "maintainer"
	License       string   "license"
	URL           string   "url"
	Size          uint64   "size"
	Sentenct      string   "sentence"
	Paragraph     string   "paragraph"
	Website       string   "website"
	Architectures []string "architectures"
	Checksum      string   "checksum"
	Category      string   "category"
}

func LoadFromFile(filename string) (*DB, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	res, err := Load(file)
	file.Close()
	return res, err
}

func Load(r io.Reader) (*DB, error) {
	decoder := json.NewDecoder(r)
	var db DB
	err := decoder.Decode(&db)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

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
