package libraries

import (
	"arduino.cc/repository/libraries/db"
	"arduino.cc/repository/libraries/metadata"
	"errors"
	git2go "github.com/libgit2/git2go"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

func CloneOrFetch(repoURL, baseFolder string) (*git2go.Repository, error) {
	parsedURL, err := url.Parse(repoURL)
	folderName := strings.NewReplacer(".git", "").Replace(parsedURL.Path)
	folderNameParts := strings.Split(folderName, "/")
	folderName = folderNameParts[len(folderNameParts)-1]
	folderName = baseFolder + "/" + folderName

	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		_, err = git2go.Clone(repoURL, folderName, &git2go.CloneOptions{})
		if err != nil {
			return nil, err
		}
	}

	repo, err := git2go.OpenRepository(folderName)
	if err != nil {
		return nil, err
	}

	origin, err := repo.LookupRemote("origin")
	err = origin.Fetch([]string{}, nil, "")

	return repo, err
}

func lastTagName(repo *git2go.Repository) (string, error) {
	referenceIterator, err := repo.NewReferenceIteratorGlob("*tags*")
	if err != nil {
		return "", err
	}

	namesIterator := referenceIterator.Names()
	var lastTagName string
	for name, err := namesIterator.Next(); err == nil; name, err = namesIterator.Next() {
		lastTagName = name
	}

	return lastTagName, nil
}

func CheckoutLastTag(repo *git2go.Repository) error {
	lastTagName, err := lastTagName(repo)
	if err != nil {
		return err
	}

	if lastTagName == "" {
		return errors.New("Repository " + repo.Workdir() + " has not tags")
	}

	ref, err := repo.LookupReference(lastTagName)
	if err != nil {
		return err
	}

	err = repo.SetHeadDetached(ref.Target(), nil, "Checking out tag "+ref.Shorthand())
	if err != nil {
		return err
	}

	return nil
}

func GenerateLibraryFromRepo(repoFolder string) (*metadata.LibraryMetadata, error) {
	bytes, err := ioutil.ReadFile(repoFolder + "library.properties")
	if err != nil {
		return nil, err
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
