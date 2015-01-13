package libraries

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"bytes"
)

func loadRepoListFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var repos []string

	reader := bufio.NewReader(file)
	var line string
	for err == nil {
		line, err = reader.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		if len(line) > 0 && line[0] != '#' {
			repos = append(repos, line)
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
	GitURLs []string
}

func (err GitURLsError) Error() string {
	error := bytes.NewBufferString("Following URL are unknown or unsupported git repos:\n")
	for _, v := range err.GitURLs {
		fmt.Fprintln(error, v)
	}

	return error.String()
}

func filterReposBy(repos []string, matcher repoMatcher) ([]string, error) {
	var filtered []string
	var wrong []string
	for _, repo := range repos {
		if matcher.Match(repo) {
			filtered = append(filtered, repo)
		} else {
			wrong = append(wrong, repo)
		}
	}
	var err error
	if len(wrong) > 0 {
		err = GitURLsError{GitURLs: wrong}
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

func toListOfUniqueRepos(repos []string) []string {
	repoSet := make(map[string]bool)

	for _, repo := range repos {
		repoSet[repo] = false
	}

	var finalRepos []string
	for key := range repoSet {
		finalRepos = append(finalRepos, key)
	}

	return finalRepos
}

func ListRepos(reposFilename string) ([]string, error) {
	repos, err := loadRepoListFromFile(reposFilename)
	if err != nil {
		return nil, err
	}

	repos, err = filterReposBy(repos, repoMatcherIfDotGit{})

	/*
		reposToVerify := filterReposBy(repos, repoMatcherIfNotDotGit{})
		githubReposToVerify := filterReposBy(reposToVerify, repoMatcherIfGithub{})
		githubOrgs, err := findGithubOrgs(githubReposToVerify)
		reposOfOrgs, err := reposFromGithubOrgs(githubOrgs)
	*/

	finalRepos := toListOfUniqueRepos(repos)

	return finalRepos, err
}
