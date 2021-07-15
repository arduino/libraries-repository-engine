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

package libraries

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/arduino/libraries-repository-engine/internal/libraries/db"

	"fmt"

	"github.com/arduino/libraries-repository-engine/internal/libraries/metadata"
	"github.com/go-git/go-git/v5"
)

// Repository represents a Git repository located on the filesystem.
type Repository struct {
	Repository *git.Repository
	FolderPath string
	URL        string
}

// CloneOrFetch returns a Repository object. If the repository is already present, it is opened. Otherwise, cloned.
func CloneOrFetch(repoMeta *Repo, folderName string) (*Repository, error) {
	repo := Repository{
		FolderPath: folderName,
		URL:        repoMeta.URL,
	}

	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		repo.Repository, err = git.PlainClone(folderName, false, &git.CloneOptions{URL: repoMeta.URL})
		if err != nil {
			return nil, err
		}
	} else {
		repo.Repository, err = git.PlainOpen(folderName)
		if err != nil {
			return nil, err
		}
	}

	tags, err := repo.Repository.Tags()
	if err != nil {
		return nil, err
	}

	for {
		tag, err := tags.Next()
		if err != nil {
			// Reached end of tags
			break
		}

		if err = repo.Repository.DeleteTag(tag.Name().Short()); err != nil {
			return nil, err
		}
	}

	if err = repo.Repository.Fetch(&git.FetchOptions{Tags: git.AllTags}); err != nil {
		return nil, err
	}

	return &repo, err
}

// GenerateLibraryFromRepo parses a repository and returns the library metadata.
func GenerateLibraryFromRepo(repo *Repository) (*metadata.LibraryMetadata, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(repo.FolderPath, "library.properties"))
	if err != nil {
		return nil, fmt.Errorf("can't read library.properties: %s", err)
	}

	library, err := metadata.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return library, nil
}

// UpdateLibrary adds a release to the library database.
func UpdateLibrary(release *db.Release, repoURL string, libraryDb *db.DB) error {
	var err error

	if !libraryDb.HasLibrary(release.LibraryName) {
		err = libraryDb.AddLibrary(&db.Library{
			Name:       release.LibraryName,
			Repository: repoURL})
		if err != nil {
			return err
		}

		err = libraryDb.Commit()
		if err != nil {
			return err
		}
	}

	if libraryDb.HasRelease(release) {
		return nil
	}

	err = libraryDb.AddRelease(release, repoURL)
	if err != nil {
		return err
	}
	err = libraryDb.Commit()
	if err != nil {
		return err
	}

	return nil
}
