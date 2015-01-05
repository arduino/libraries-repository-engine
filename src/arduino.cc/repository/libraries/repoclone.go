package libraries

import (
	git2go "github.com/libgit2/git2go"
	"net/url"
	"strings"
)

// TODO complete implementation
func CloneOrPull(repoURL, baseFolder string) error {
	parsed, err := url.Parse(repoURL)
	folderName := strings.Split(strings.Split(parsed.Path, "/")[2], ".")[0]

	_, err = git2go.Clone(repoURL, baseFolder+"/"+folderName, &git2go.CloneOptions{})

	return err
}
