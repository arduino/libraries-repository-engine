package libraries

import (
	git2go "github.com/libgit2/git2go"
	"net/url"
	"strings"
	"os"
)

func CloneOrFetch(repoURL, baseFolder string) (*git2go.Repository, error) {
	parsedURL, err := url.Parse(repoURL)
	folderName := strings.NewReplacer(".git", "").Replace(parsedURL.Path)
	folderNameParts := strings.Split(folderName, "/")
	folderName = folderNameParts[len(folderNameParts)-1]
	folderName = baseFolder+"/"+folderName

	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		_, err = git2go.Clone(repoURL, folderName, &git2go.CloneOptions{})
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
