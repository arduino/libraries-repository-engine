// This file is part of libraries-repository-engine.
//
// Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

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

// LoadRepoListFromFile returns an unfiltered list of library registry entries loaded from the given data file.
func LoadRepoListFromFile(filename string) ([]*Repo, error) {
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
				URL:         url,
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

func (repoMatcherIfDotGit) Match(url string) bool {
	return strings.Index(url, "https://") == 0 && strings.LastIndex(url, ".git") == len(url)-len(".git")
}

// GitURLsError is the type for the unknown or unsupported repositories data.
type GitURLsError struct {
	Repos []*Repo
}

// Repo is the type for the library repository data.
type Repo struct {
	URL         string
	Types       []string
	LibraryName string
}

// AsFolder returns the URL of the repo as path, without protocol prefix or suffix.
// For example if the repo URL is https://github.com/example/lib.git this function
// will return "github.com/example/lib"
func (repo *Repo) AsFolder() (string, error) {
	u, err := url.Parse(repo.URL)
	if err != nil {
		return "", err
	}
	folderName := strings.Replace(u.Path, ".git", "", -1)
	folderName = filepath.Join(u.Host, folderName)
	return folderName, nil
}

// ReposByURL is the type for the libraries repository data.
type ReposByURL []*Repo

func (r ReposByURL) Len() int {
	return len(r)
}

func (r ReposByURL) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ReposByURL) Less(i, j int) bool {
	return r[i].URL < r[j].URL
}

func (err GitURLsError) Error() string {
	error := bytes.NewBufferString("Following URL are unknown or unsupported git repos:\n")
	for _, v := range err.Repos {
		fmt.Fprintln(error, v.URL)
	}

	return error.String()
}

func filterReposBy(repos []*Repo, matcher repoMatcher) ([]*Repo, error) {
	var filtered []*Repo
	var wrong []*Repo
	for _, repo := range repos {
		if matcher.Match(repo.URL) {
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

func toListOfUniqueRepos(repos []*Repo) []*Repo {
	repoMap := make(map[string]*Repo)
	var finalRepos []*Repo

	for _, repo := range repos {
		if _, contains := repoMap[repo.URL]; !contains {
			finalRepos = append(finalRepos, repo)
			repoMap[repo.URL] = repo
		}
	}

	return finalRepos
}

// ListRepos returns a filtered list of library registry entries loaded from the given data file.
func ListRepos(reposFilename string) ([]*Repo, error) {
	repos, err := LoadRepoListFromFile(reposFilename)
	if err != nil {
		return nil, err
	}

	repos, err = filterReposBy(repos, repoMatcherIfDotGit{})

	finalRepos := toListOfUniqueRepos(repos)

	return finalRepos, err
}
