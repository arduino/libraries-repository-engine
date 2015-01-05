package libraries

import (
	"os"
	"bufio"
	"strings"
	"github.com/cmaglie/go-github/github"
	"net/url"
)

type Repo struct {
	GitURL string
}

func loadRepoListFromFile(filename string) ([]Repo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var repos []Repo

	reader := bufio.NewReader(file)
	var line string
	for err == nil {
		line, err = reader.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		if len(line) > 0 && line[0] != '#' {
			repos = append(repos, Repo{GitURL: line})
		}

	}

	return repos, nil
}

type repoMatcher interface {
	Match(Repo) bool
}

type repoMatcherIfDotGit struct {}

func (_ repoMatcherIfDotGit) Match(r Repo) bool {
	return strings.LastIndex(r.GitURL, ".git") == len(r.GitURL)-len(".git")
}

type repoMatcherIfNotDotGit struct {}

func (_ repoMatcherIfNotDotGit) Match(r Repo) bool {
	return !repoMatcherIfDotGit{}.Match(r)
}

type repoMatcherIfGithub struct {}

func (_ repoMatcherIfGithub) Match(r Repo) bool {
	return strings.Index(r.GitURL, "//github.com") != -1 || strings.Index(r.GitURL, "@github.com") != -1
}

func filterReposBy(repos []Repo, matcher repoMatcher) []Repo {
	var filtered []Repo
	for _, repo := range repos {
		if matcher.Match(repo) {
			filtered = append(filtered, repo)
		}
	}
	return filtered
}

func reposFromGithubOrgs(orgs []*github.Organization) ([]Repo, error) {
	client := github.NewClient(nil)
	var repos []Repo
	for _, org := range orgs {
		repositories, _, err := client.Repositories.ListByOrg(*org.Login, &github.RepositoryListByOrgOptions{})
		if err != nil {
			return nil, err
		}
		for _, repository := range repositories {
			repos = append(repos, Repo{GitURL: *repository.SSHURL})
		}
	}

	return repos, nil
}

func findGithubOrgs(repos []Repo) (orgs []*github.Organization, err error) {
	client := github.NewClient(nil)
	for _, repo := range repos {
		url, err := url.Parse(repo.GitURL)
		if err != nil {
			return nil, err
		}
		orgName := strings.Split(url.Path, "/")[1]
		org, _, err := client.Organizations.Get(orgName)
		if err == nil {
			orgs = append(orgs, org)
		}
	}
	return orgs, nil
}

func ListRepos(reposFilename string) ([]string, error) {
	repos, err := loadRepoListFromFile(reposFilename)
	if err != nil {
		return nil, err
	}

	reposAlreadyOk := filterReposBy(repos, repoMatcherIfDotGit{})

	reposToVerify := filterReposBy(repos, repoMatcherIfNotDotGit{})
	githubReposToVerify := filterReposBy(reposToVerify, repoMatcherIfGithub{})
	githubOrgs, err := findGithubOrgs(githubReposToVerify)
	reposOfOrgs, err := reposFromGithubOrgs(githubOrgs)

	repoSet := make(map[string]bool)

	loadSet := func(m map[string]bool, repos []Repo) {
		for _, repo := range repos {
			m[repo.GitURL] = false
		}
	}
	loadSet(repoSet, reposAlreadyOk)
	loadSet(repoSet, reposOfOrgs)

	var finalRepos []string
	for key := range repoSet {
		finalRepos = append(finalRepos, key)
	}

	return finalRepos, err
}

