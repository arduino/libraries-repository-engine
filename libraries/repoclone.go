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
	"github.com/go-git/go-git/v5"
)

// Repository represents a Git repository located on the filesystem.
type Repository struct {
	Repository *git.Repository
	FolderPath string
	URL        string
}

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

		if err = repo.Repository.DeleteTag(strings.TrimPrefix(tag.Name().String(), "refs/tags/")); err != nil {
			return nil, err
		}
	}

	if err = repo.Repository.Fetch(&git.FetchOptions{Tags: git.AllTags}); err != nil {
		return nil, err
	}

	return &repo, err
}

func GenerateLibraryFromRepo(repo *Repository) (*metadata.LibraryMetadata, error) {
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
