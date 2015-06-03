package libraries

import (
	"arduino.cc/repository/libraries/db"
	"arduino.cc/repository/libraries/git"
	"arduino.cc/repository/libraries/metadata"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
)

func CloneOrFetch(repoURL, baseFolder string) (string, error) {
	parsedURL, err := url.Parse(repoURL)
	folderName := strings.NewReplacer(".git", "").Replace(parsedURL.Path)
	folderNameParts := strings.Split(folderName, "/")
	folderName = folderNameParts[len(folderNameParts)-1]
	folderName = path.Join(baseFolder, folderName)

	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = git.Clone(repoURL, folderName)
		if err != nil {
			return "", err
		}
	}

	tags, err := git.ListTags(folderName)
	if err != nil {
		return "", err
	}
	for _, tag := range tags {
		err = git.RemoveTag(folderName, tag)
		if err != nil {
			return "", err
		}
	}

	err = git.Fetch(folderName)
	if err != nil {
		return "", err
	}

	return folderName, err
}

func ListTags(folderName string) ([]string, error) {
	return git.ListTags(folderName)
}

func lastTagName(folderName string) (string, error) {
	return git.LastTag(folderName)
}

func CheckoutTag(folderName string, tagName string) error {
	return git.CheckoutTag(folderName, tagName)
}

func CheckoutLastTag(folderName string) error {
	lastTagName, err := lastTagName(folderName)
	if err != nil {
		return err
	}

	if lastTagName == "" {
		return errors.New("Repository " + folderName + " has not tags")
	}

	return git.CheckoutTag(folderName, lastTagName)
}

func GenerateLibraryFromRepo(repoFolder string, repo *Repo) (*metadata.LibraryMetadata, error) {
	bytes, err := ioutil.ReadFile(path.Join(repoFolder, "library.properties"))
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
