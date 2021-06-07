// This file is part of libraries-repository-engine.
//
// Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package db

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

// DB is the libraries database
type DB struct {
	Libraries []*Library
	Releases  []*Release

	libraryFile string
	mutex       sync.Mutex
}

// Library is an Arduino library
type Library struct {
	Name         string
	Repository   string
	SupportLevel string

	// Category of the latest release of the library
	LatestCategory string
}

// Release is a library release
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
	Includes        []string
	Dependencies    []*Dependency
	Log             string
}

// Dependency is a library dependency
type Dependency struct {
	Name    string
	Version string
}

// New returns a new DB object.
func New(libraryFile string) *DB {
	return &DB{libraryFile: libraryFile}
}

// AddLibrary adds a library to the database.
func (db *DB) AddLibrary(library *Library) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	found, _ := db.findLibrary(library.Name)
	if found != nil {
		return errors.New("library already exists")
	}
	db.Libraries = append(db.Libraries, library)
	return nil
}

// HasLibrary returns whether the database already contains the given library.
func (db *DB) HasLibrary(libraryName string) bool {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.hasLibrary(libraryName)
}

func (db *DB) hasLibrary(libraryName string) bool {
	found, _ := db.findLibrary(libraryName)
	return found != nil
}

// FindLibrary returns the Library object for the given name.
func (db *DB) FindLibrary(libraryName string) (*Library, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.findLibrary(libraryName)
}

func (db *DB) findLibrary(libraryName string) (*Library, error) {
	for _, lib := range db.Libraries {
		if lib.Name == libraryName {
			return lib, nil
		}
	}
	return nil, errors.New("library not found")
}

// AddRelease adds a library release to the database.
func (db *DB) AddRelease(release *Release, repoURL string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	lib, err := db.findLibrary(release.LibraryName)
	if err != nil {
		return err
	}
	lib.Repository = repoURL

	if db.hasRelease(release) {
		return errors.New("release already exists")
	}
	db.Releases = append(db.Releases, release)

	// Update LatestCategory with the Category of the latest release
	last, err := db.findLatestReleaseOfLibrary(lib)
	if err != nil {
		return err
	}
	lib.LatestCategory = last.Category

	return nil
}

// HasReleaseByNameVersion returns whether the database contains a release for the given library and version number.
func (db *DB) HasReleaseByNameVersion(libraryName string, libraryVersion string) bool {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.hasReleaseByNameVersion(libraryName, libraryVersion)
}

func (db *DB) hasReleaseByNameVersion(libraryName string, libraryVersion string) bool {
	found, _ := db.findReleaseByNameVersion(libraryName, libraryVersion)
	return found != nil
}

// HasRelease returns whether the database already contains the given Release object.
func (db *DB) HasRelease(release *Release) bool {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.hasRelease(release)
}

func (db *DB) hasRelease(release *Release) bool {
	return db.hasReleaseByNameVersion(release.LibraryName, release.Version.String())
}

// FindRelease returns the Release object from the database that matches the given object.
func (db *DB) FindRelease(release *Release) (*Release, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.findReleaseByNameVersion(release.LibraryName, release.Version.String())
}

func (db *DB) findReleaseByNameVersion(libraryName string, libraryVersion string) (*Release, error) {
	for _, r := range db.Releases {
		if r.LibraryName == libraryName && r.Version.String() == libraryVersion {
			return r, nil
		}
	}
	return nil, errors.New("library not found")
}

// LoadFromFile returns a DB object loaded from the given filename.
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

// Load returns a DB object loaded from the given reader.
func Load(r io.Reader) (*DB, error) {
	decoder := json.NewDecoder(r)
	db := new(DB)
	err := decoder.Decode(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// SaveToFile saves the database to a file.
func (db *DB) SaveToFile() error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	file, err := os.Create(db.libraryFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return db.save(file)
}

// Save writes the database via the given writer.
func (db *DB) Save(r io.Writer) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.save(r)
}

func (db *DB) save(r io.Writer) error {
	buff, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	_, err = r.Write(buff)
	return err
}

func (db *DB) findLatestReleaseOfLibrary(lib *Library) (*Release, error) {
	var found *Release
	for _, rel := range db.findReleasesOfLibrary(lib) {
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

// FindReleasesOfLibrary returns the database's releases for the given Library object.
func (db *DB) FindReleasesOfLibrary(lib *Library) []*Release {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	return db.findReleasesOfLibrary(lib)
}

func (db *DB) findReleasesOfLibrary(lib *Library) []*Release {
	var releases []*Release
	for _, rel := range db.Releases {
		if rel.LibraryName != lib.Name {
			continue
		}
		releases = append(releases, rel)
	}
	return releases
}

// Commit saves the database to disk.
func (db *DB) Commit() error {
	return db.SaveToFile()
}

// Init loads a database from file and returns it.
func Init(libraryFile string) *DB {
	libs, err := LoadFromFile(libraryFile)
	if err != nil {
		log.Print(err)
		log.Print("starting with an empty DB")
		return New(libraryFile)
	}
	log.Printf("Loaded %v libraries from DB", len(libs.Libraries))
	return libs
}
