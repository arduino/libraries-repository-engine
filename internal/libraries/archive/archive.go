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

// Package archive handles the library release archive.
package archive

import (
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/arduino/libraries-repository-engine/internal/configuration"
	"github.com/arduino/libraries-repository-engine/internal/libraries"
	"github.com/arduino/libraries-repository-engine/internal/libraries/hash"
	"github.com/arduino/libraries-repository-engine/internal/libraries/metadata"
	"github.com/arduino/libraries-repository-engine/internal/libraries/zip"
)

// Archive is the type for library release archive data.
type Archive struct {
	SourcePath string
	RootName   string // Name of the root folder inside the archive.
	FileName   string
	Path       string // Full path of the archive.
	URL        string // URL the archive will have on the download server.
	Size       int64
	Checksum   string
}

// New initializes and returns an Archive object.
func New(repository *libraries.Repository, libraryMetadata *metadata.LibraryMetadata, config *configuration.Config) (*Archive, error) {
	repositoryURLData, err := url.Parse(repository.URL)
	if err != nil {
		return nil, err
	}
	// e.g., https://github.com/arduino-libraries/Servo.git -> github.com
	repositoryHost := repositoryURLData.Host
	// e.g., https://github.com/arduino-libraries/Servo.git -> arduino-libraries
	repositoryParent := filepath.Base(filepath.Dir(repositoryURLData.Path))

	// Unlike the other path components, the filename is based on library name, not repository name URL.
	fileName := zipFolderName(libraryMetadata) + ".zip"

	return &Archive{
		SourcePath: repository.FolderPath,
		RootName:   zipFolderName(libraryMetadata),
		FileName:   fileName,
		Path:       filepath.Join(config.LibrariesFolder, repositoryHost, repositoryParent, fileName),
		URL:        config.BaseDownloadURL + repositoryHost + "/" + repositoryParent + "/" + fileName,
	}, nil
}

// Create makes an archive file according to the data of the Archive object and updates the object with the size and
// checksum for the resulting file.
func (archive *Archive) Create() error {
	err := os.MkdirAll(filepath.Dir(archive.Path), os.FileMode(0755))
	if err != nil {
		return err
	}

	if err := zip.Directory(archive.SourcePath, archive.RootName, archive.Path); err != nil {
		os.Remove(archive.Path)
		return err
	}

	size, checksum, err := getSizeAndCalculateChecksum(archive.Path)
	if err != nil {
		return err
	}
	archive.Size = size
	archive.Checksum = checksum

	return nil
}

var zipFolderNamePattern = regexp.MustCompile("[^a-zA-Z0-9]")

// zipFolderName returns the name to use for the folder.
func zipFolderName(library *metadata.LibraryMetadata) string {
	return zipFolderNamePattern.ReplaceAllString(library.Name, "_") + "-" + library.Version
}

// getSizeAndCalculateChecksum returns the size and SHA-256 checksum for the given file.
func getSizeAndCalculateChecksum(filePath string) (int64, string, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return -1, "", err
	}

	size := info.Size()

	checksum, err := hash.Checksum(filePath)
	if err != nil {
		return -1, "", err
	}

	return size, checksum, nil
}
