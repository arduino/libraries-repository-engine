package libraries

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func loadRepoListFromFile(filename string) ([]*Repo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var repos []*Repo

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] != '#' {
			split := strings.Split(line, "|")
			url := split[0]
			types := strings.Split(split[1], ",")
			name := split[2]
			repos = append(repos, &Repo{
				Url:         url,
				Types:       types,
				LibraryName: name,
			})
		}
	}

	return repos, nil
}

type repoMatcher interface {
	Match(string) bool
}

type repoMatcherIfDotGit struct{}

func (_ repoMatcherIfDotGit) Match(url string) bool {
	return strings.Index(url, "https://") == 0 && strings.LastIndex(url, ".git") == len(url)-len(".git")
}

/*
type repoMatcherIfNotDotGit struct{}

func (_ repoMatcherIfNotDotGit) Match(r string) bool {
	return !repoMatcherIfDotGit{}.Match(r)
}

type repoMatcherIfGithub struct{}

func (_ repoMatcherIfGithub) Match(r string) bool {
	return strings.Index(r, "//github.com") != -1 || strings.Index(r, "@github.com") != -1
}
*/

type GitURLsError struct {
	Repos []*Repo
}

type Repo struct {
	Url         string
	Types       []string
	LibraryName string
}

// AsFolder returns the URL of the repo as path, without protocol prefix or suffix.
// For example if the repo URL is https://github.com/example/lib.git this function
// will return "github.com/example/lib"
func (repo *Repo) AsFolder() (string, error) {
	u, err := url.Parse(repo.Url)
	if err != nil {
		return "", err
	}
	folderName := strings.Replace(u.Path, ".git", "", -1)
	folderName = filepath.Join(u.Host, folderName)
	return folderName, nil
}

type ReposByUrl []*Repo

func (r ReposByUrl) Len() int {
	return len(r)
}

func (r ReposByUrl) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ReposByUrl) Less(i, j int) bool {
	return r[i].Url < r[j].Url
}

func (err GitURLsError) Error() string {
	error := bytes.NewBufferString("Following URL are unknown or unsupported git repos:\n")
	for _, v := range err.Repos {
		fmt.Fprintln(error, v.Url)
	}

	return error.String()
}

func filterReposBy(repos []*Repo, matcher repoMatcher) ([]*Repo, error) {
	var filtered []*Repo
	var wrong []*Repo
	for _, repo := range repos {
		if matcher.Match(repo.Url) {
			filtered = append(filtered, repo)
		} else {
			wrong = append(wrong, repo)
		}
	}
	var err error
	if len(wrong) > 0 {
		err = GitURLsError{wrong}
	}
	return filtered, err
}

/*
func newGithubClient() *github.Client {
	gh_auth := &oauth.Transport{Token: &oauth.Token{AccessToken: config.GithubAuthToken()}}
	return github.NewClient(gh_auth.Client())
}

func reposFromGithubOrgs(orgs []*github.Organization) ([]string, error) {
	client := newGithubClient()
	var repos []string
	for _, org := range orgs {
		repositories, _, err := client.Repositories.ListByOrg(*org.Login, &github.RepositoryListByOrgOptions{})
		if err != nil {
			return nil, err
		}
		for _, repository := range repositories {
			repos = append(repos, *repository.CloneURL)
		}
	}

	return repos, nil
}

func findGithubOrgs(repos []string) (orgs []*github.Organization, err error) {
	client := newGithubClient()
	for _, repo := range repos {
		parsedURL, err := url.Parse(repo)
		if err != nil {
			return nil, err
		}
		orgName := strings.Split(parsedURL.Path, "/")[1]
		org, _, err := client.Organizations.Get(orgName)
		if err == nil {
			orgs = append(orgs, org)
		}
	}
	return orgs, nil
}
*/

func toListOfUniqueRepos(repos []*Repo) []*Repo {
	repoMap := make(map[string]*Repo)
	var finalRepos []*Repo

	for _, repo := range repos {
		if _, contains := repoMap[repo.Url]; !contains {
			finalRepos = append(finalRepos, repo)
			repoMap[repo.Url] = repo
		}
	}

	return finalRepos
}

func ListRepos(reposFilename string) ([]*Repo, error) {
	repos, err := loadRepoListFromFile(reposFilename)
	if err != nil {
		return nil, err
	}

	repos, err = filterReposBy(repos, repoMatcherIfDotGit{})

	finalRepos := toListOfUniqueRepos(repos)

	return finalRepos, err
}
