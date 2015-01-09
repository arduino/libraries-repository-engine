package db

import "encoding/json"
import "errors"
import "io"
import "os"

// The libraries DB
type DB struct {
	Libraries []*Library
	Releases  []*Release
}

// A library
type Library struct {
	Name       string
	Repository string
}

// A release
type Release struct {
	LibraryName   string // The library name
	Version       Version
	Author        string
	Maintainer    string
	License       string
	Sentence      string
	Paragraph     string
	Website       string
	Category      string
	Architectures []string

	URL             string
	ArchiveFileName string
	Size            int64
	Checksum        string
}

func New() *DB {
	return &DB{}
}

func (db *DB) AddLibrary(library *Library) error {
	found, _ := db.FindLibrary(library.Name)
	if found != nil {
		return errors.New("library alredy existent")
	}
	db.Libraries = append(db.Libraries, library)
	return nil
}

func (db *DB) HasLibrary(libraryName string) bool {
	found, _ := db.FindLibrary(libraryName)
	return found != nil
}

func (db *DB) FindLibrary(libraryName string) (*Library, error) {
	for _, lib := range db.Libraries {
		if lib.Name == libraryName {
			return lib, nil
		}
	}
	return nil, errors.New("library not found")
}

func (db *DB) AddRelease(release *Release) error {
	if !db.HasLibrary(release.LibraryName) {
		return errors.New("released library not found")
	}
	if db.HasRelease(*release) {
		return errors.New("release already existent")
	}
	db.Releases = append(db.Releases, release)
	return nil
}

func (db *DB) HasRelease(release Release) bool {
	found, _ := db.FindRelease(release)
	return found != nil
}

func (db *DB) FindRelease(release Release) (*Release, error) {
	for _, r := range db.Releases {
		if r.LibraryName == release.LibraryName && r.Version == release.Version {
			return r, nil
		}
	}
	return nil, errors.New("library not found")
}

func LoadFromFile(filename string) (*DB, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Load(file)
}

func Load(r io.Reader) (*DB, error) {
	decoder := json.NewDecoder(r)
	db := new(DB)
	err := decoder.Decode(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return db.Save(file)
}

func (db *DB) Save(r io.Writer) error {
	buff, err := json.MarshalIndent(*db, "", "  ")
	if err != nil {
		return err
	}
	_, err = r.Write(buff)
	return err
}

func (db *DB) FindLatestReleaseOfLibrary(lib *Library) (*Release, error) {
	var found *Release = nil
	for _, rel := range db.Releases {
		if rel.LibraryName != lib.Name {
			continue
		}
		if found == nil {
			found = rel
		} else {
			if less, err := found.Version.Less(rel.Version); err != nil {
				return nil, err
			} else if less {
				found = rel
			}
		}
	}
	return found, nil
}

// vi:ts=2
