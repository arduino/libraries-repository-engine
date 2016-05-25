package libraries

import (
	"arduino.cc/repository/libraries/db"
	"arduino.cc/repository/libraries/metadata"
	"errors"
	"github.com/arduino/arduino-modules/git"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func CloneOrFetch(repoURL, baseFolder string) (*git.Repository, error) {
	parsedURL, err := url.Parse(repoURL)
	folderName := strings.NewReplacer(".git", "").Replace(parsedURL.Path)
	folderNameParts := strings.Split(folderName, "/")[1:]
	folderName = filepath.Join(baseFolder, filepath.Join(folderNameParts...))

	var repo *git.Repository
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		repo, err = git.Clone(repoURL, folderName)
		if err != nil {
			return nil, err
		}
	} else {
		repo = &git.Repository{FolderPath: folderName}
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

func CheckoutLastTag(repo *git.Repository) error {
	tags, err := repo.ListTags()
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return errors.New("No tags in repository " + repo.FolderPath)
	}

	lastTagName := tags[len(tags)-1]
	if lastTagName == "" {
		return errors.New("Repository " + repo.FolderPath + " has not tags")
	}

	return repo.CheckoutTag(lastTagName)
}

func GenerateLibraryFromRepo(repoFolder string, repo *Repo) (*metadata.LibraryMetadata, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(repoFolder, "library.properties"))
	if err != nil {
		return nil, err
	}

	library, err := metadata.Parse(bytes)
	if err != nil {
		return nil, err
	}
	library.Types = repo.Types

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

func UpdateLibrary(release *db.Release, libraryDb *db.DB) error {
	var err error

	if !libraryDb.HasLibrary(release.LibraryName) {
		err = libraryDb.AddLibrary(&db.Library{Name: release.LibraryName})
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

	err = libraryDb.AddRelease(release)
	if err != nil {
		return err
	}
	err = libraryDb.Commit()
	if err != nil {
		return err
	}

	return nil
}
