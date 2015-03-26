package libraries

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func loadRepoListFromFile(filename string) ([]*Repo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var repos []*Repo

	reader := bufio.NewReader(file)
	var line string
	for err == nil {
		line, err = reader.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		if len(line) > 0 && line[0] != '#' {
			lineParts := strings.Split(line, "\t")
			url := lineParts[0]
			types := strings.Split(lineParts[1], ",")
			repos = append(repos, &Repo{url, types})
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
	Url   string
	Types []string
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
	repoSet := make(map[string]*Repo)

	for _, repo := range repos {
		repoSet[repo.Url] = repo
	}

	var finalRepos []*Repo
	for _, value := range repoSet {
		finalRepos = append(finalRepos, value)
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
