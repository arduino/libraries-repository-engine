package db

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

// The libraries DB
type DB struct {
	libraryFile string
	Libraries   []*Library
	Releases    []*Release
}

// A library
type Library struct {
	Name         string
	Repository   string
	SupportLevel string

	// Category of the latest release of the library
	LatestCategory string
}

// A release
type Release struct {
	LibraryName     string // The library name
	Version         Version
	Author          string
	Maintainer      string
	License         string
	Sentence        string
	Paragraph       string
	Website         string
	Category        string
	Architectures   []string
	Types           []string
	URL             string
	ArchiveFileName string
	Size            int64
	Checksum        string
}

func New(libraryFile string) *DB {
	return &DB{libraryFile: libraryFile}
}

func (db *DB) AddLibrary(library *Library) error {
	found, _ := db.FindLibrary(library.Name)
	if found != nil {
		return errors.New("library already exists")
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
	lib, err := db.FindLibrary(release.LibraryName)
	if err != nil {
		return err
	}

	if db.HasRelease(release) {
		return errors.New("release already exists")
	}
	db.Releases = append(db.Releases, release)

	// Update LatestCategory with the Category of the latest release
	last, err := db.FindLatestReleaseOfLibrary(lib)
	if err != nil {
		return err
	}
	lib.LatestCategory = last.Category

	return nil
}

func (db *DB) HasReleaseByNameVersion(libraryName string, libraryVersion string) bool {
	found, _ := db.FindReleaseByNameVersion(libraryName, libraryVersion)
	return found != nil
}

func (db *DB) HasRelease(release *Release) bool {
	return db.HasReleaseByNameVersion(release.LibraryName, release.Version.String())
}

func (db *DB) FindRelease(release *Release) (*Release, error) {
	return db.FindReleaseByNameVersion(release.LibraryName, release.Version.String())
}

func (db *DB) FindReleaseByNameVersion(libraryName string, libraryVersion string) (*Release, error) {
	for _, r := range db.Releases {
		if r.LibraryName == libraryName && r.Version.String() == libraryVersion {
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
	db, err := Load(file)
	if err != nil {
		return nil, err
	}
	db.libraryFile = filename
	return db, nil
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

func (db *DB) SaveToFile() error {
	file, err := os.Create(db.libraryFile)
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
	for _, rel := range db.FindReleasesOfLibrary(lib) {
		if found == nil {
			found = rel
			continue
		}
		if less, err := found.Version.Less(rel.Version); err != nil {
			return nil, err
		} else if less {
			found = rel
		}
	}
	return found, nil
}

func (db *DB) FindReleasesOfLibrary(lib *Library) []*Release {
	var releases []*Release
	for _, rel := range db.Releases {
		if rel.LibraryName != lib.Name {
			continue
		}
		releases = append(releases, rel)
	}
	return releases
}

func (db *DB) Commit() error {
	return db.SaveToFile()
}

func Init(libraryFile string) *DB {
	if libs, err := LoadFromFile(libraryFile); err != nil {
		log.Print(err)
		log.Print("starting with an empty DB")
		return New(libraryFile)
	} else {
		log.Printf("Loaded %v libraries from DB", len(libs.Libraries))
		return libs
	}
}
