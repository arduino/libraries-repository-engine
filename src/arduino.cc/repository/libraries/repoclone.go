package libraries

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"arduino.cc/repository/libraries/db"

	"fmt"

	"arduino.cc/repository/libraries/metadata"
	"github.com/arduino/arduino-modules/git"
)

func CloneOrFetch(repoMeta *Repo, folderName string) (*git.Repository, error) {
	var repo *git.Repository
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		repo, err = git.Clone(repoMeta.Url, folderName)
		if err != nil {
			return nil, err
		}
	} else {
		repo = &git.Repository{FolderPath: folderName, URL: repoMeta.Url}
	}

	tags, err := repo.ListTags()
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		if err = repo.RemoveTag(tag); err != nil {
			return nil, err
		}
	}

	if err = repo.Fetch(); err != nil {
		return nil, err
	}

	return repo, err
}

func GenerateLibraryFromRepo(repo *git.Repository) (*metadata.LibraryMetadata, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(repo.FolderPath, "library.properties"))
	if err != nil {
		return nil, fmt.Errorf("can't read library.properties: %s", err)
	}

	library, err := metadata.Parse(bytes)
	if err != nil {
		return nil, err
	}

	libraryErrors := library.Validate()
	if len(libraryErrors) > 0 {
		var errorsString []string
		for _, error := range libraryErrors {
			errorsString = append(errorsString, error.Error())
		}
		combinedErrors := strings.Join(errorsString, ",")
		return nil, errors.New(combinedErrors)
	}

	return library, nil
}

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
