package main

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	//"os/exec"
	"bytes"
	"encoding/base64"
	"github.com/vaughan0/go-ini"
	s "strings"
)

// Configuration part
// TODO: Move into his own file
var gh_auth_token = "e282d0ab13c38a4303e65620aeab13c2beba3385"
var org = "arlibs"
var gh_callbackUrl = "http://kungfu.bug.st:8088/github/event/"
var localGitFolder = "git/"

// Global github client
var gh_auth = &oauth.Transport{Token: &oauth.Token{AccessToken: gh_auth_token}}
var gh = github.NewClient(gh_auth.Client())

// Github event web hook.
func GithubEventHook(c *gin.Context) {
	eventType := c.Request.Header.Get("X-GitHub-Event")

	switch eventType {
	case "ping":
		// Ping event has only "zen" and "hook_id" values
		var ping struct {
			Zen    *string `json:"zen"`
			HookID *int    `json:"hook_id"`
		}
		c.Bind(&ping)
		c.String(200, "Received ping from github: "+*ping.Zen)
		return

	case "pull_request":
		var event github.PullRequestEvent
		c.Bind(&event)
		CheckRelease(c, event.PullRequest)
		return
	}

	c.String(200, "Received "+eventType+" from github. Ignoring...")
}

func CheckRelease(c *gin.Context, pull *github.PullRequest) {
	//commits := *pull.Commits
	title := *pull.Title

	// Is a release request?
	if s.HasPrefix(title, "[RELEASE] ") {

		// Check if the release number is the same inside library
		version := title[10:]
		fmt.Println("Doing release!")
		fmt.Println("  release in PR: '" + version + "'")

		fmt.Println(title)

		head := *pull.Head
		headRepo := *head.Repo
		fmt.Println("  release sha  : " + *head.SHA)
		fmt.Println("  repository   : " + *head.Repo.FullName)

		// Get library.properties from pull request HEAD
		getContentOpts := &github.RepositoryContentGetOptions{
			Ref: *head.SHA,
		}
		libPropContent, _, _, err := gh.Repositories.GetContents(*headRepo.Owner.Login, *headRepo.Name, "library.properties", getContentOpts)
		if err != nil || libPropContent == nil {
			c.JSON(500, gin.H{
				"result":   "error",
				"message":  "cannot fetch library.properties",
				"gh_error": err,
			})
			return
		}

		libPropertiesData, err := base64.StdEncoding.DecodeString(*libPropContent.Content)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		reader := bytes.NewReader(libPropertiesData)
		properties, err := ini.Load(reader)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		libName, _ := properties.Get("", "name")
		fmt.Printf("  library.properties contains: %q\n", libName)

		resultMsg := "Thanks!"
		/*
			_, _, err := CommentOnPullRequest(pull, resultMsg)
			if err != nil {
				fmt.Println(github.Stringify(err))
				c.JSON(500, gin.H{
					"result":   "error",
					"message":  "error sending comment message",
					"gh_error": err,
				})
				return
			}
		*/
		fmt.Println(resultMsg)
		c.String(200, "Received pull_request from github.")
	}
}

func CommentOnPullRequest(pull *github.PullRequest, text string) (*github.IssueComment, *github.Response, error) {
	comment := &github.IssueComment{
		Body: github.String(text),
	}
	owner := *pull.Base.Repo.Owner.Login
	repo := *pull.Base.Repo.Name
	number := *pull.Number
	return gh.Issues.CreateComment(owner, repo, number, comment)
}

// Create a new repository for the specified library.
// A new empty repository is created on the library manager organization
// and it is connected via web hook to the library manager.
func CreateLibrary(c *gin.Context) {
	name := c.Params.ByName("name")

	// Create a new repository named as the requested library
	repository := &github.Repository{
		Name:      github.String(name),
		HasIssues: github.Bool(true),
		AutoInit:  github.Bool(true), // To be removed in favor of our custom PushInitialEmptyRepository
	}
	newRepository, _, err := gh.Repositories.Create(org, repository)
	if err != nil {
		c.JSON(500, gin.H{
			"result":   "error",
			"message":  "error creating repository",
			"gh_error": err,
		})
		return
	}

	// Activate web hook for ALL events on this repository
	hook := &github.Hook{
		// must be "web" for web hooks
		Name:   github.String("web"),
		Events: []string{"*"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          gh_callbackUrl + name,
			"content_type": "json",
		},
	}
	newHook, _, err := gh.Repositories.CreateHook(org, name, hook)
	if err != nil {
		c.JSON(500, gin.H{
			"result":     "error",
			"message":    "error creating repository web hook",
			"repository": newRepository,
			"gh_error":   err,
		})
		return
	}

	// Push initial state on the repository
	//PushInitialEmptyRepository(c, newRepository)

	c.JSON(200, gin.H{
		"result":     "ok",
		"message":    "created repository " + name,
		"repository": newRepository,
		"hook":       newHook,
	})
}

// Push the inital empty repository with readme for developers
func PushInitialEmptyRepository(c *gin.Context, repo *github.Repository) {
	Run(localGitFolder, "mkdir", *repo.Name)
	gitFolder := localGitFolder + *repo.Name
	Run(gitFolder, "git", "init")
	Run(gitFolder, "touch", "README.md")
	Run(gitFolder, "git", "add", "README.md")
	Run(gitFolder, "git", "commit", "-m", "Initialized library repostiory")
	Run(gitFolder, "git", "remote", "add", "origin", "git@github-as-arlib0:arlibs/"+*repo.Name+".git")
	Run(gitFolder, "git", "push", "-u", "origin", "master")
}

func Run(workdir string, name string, arg ...string) int {
/*
	cmd := exec.Command(name, arg)
	cmd.Dir = workdir
	cmd.Start()
	return cmd.Wait()
*/
  return 0
}

func ListAllLibraries(c *gin.Context) {
	teams, _, err := gh.Organizations.ListTeams("arlibs", nil)
	if err != nil {
		c.JSON(500, gin.H{
			"result":   "error",
			"message":  "could not fetch organization teams",
			"gh_error": err,
		})
		return
	}

	team := teams[0] // The only team available should be "owners"
	// fmt.Println("Teams : ", *team.Name)

	repos, _, err := gh.Organizations.ListTeamRepos(*team.ID, nil)
	if err != nil {
		c.JSON(500, gin.H{
			"result":   "error",
			"message":  "could not get organization repositories",
			"gh_error": err,
		})
		return
	}
	libraries := make([]struct {
		Name   string
		GitURL string
	}, len(repos))
	for i, repo := range repos {
		libraries[i].Name = *repo.Name
		libraries[i].GitURL = *repo.GitURL
	}

	c.JSON(200, gin.H{
		"result":    "ok",
		"libraries": libraries,
	})
}

func main() {
	r := gin.Default()

	r.GET("/libraries/list/", ListAllLibraries)
	r.GET("/libraries/create/:name", CreateLibrary)
	r.POST("/github/event/:name", GithubEventHook)

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")

	org, _, err := gh.Organizations.Get("arlibs")
	if err != nil {
		fmt.Printf("error: %v\n\n", err)
	} else {
		//fmt.Printf("%v\n\n", github.Stringify(org))
	}

	fmt.Println("Organization: ", *org.Login)
	fmt.Println("Repositories: ", github.Stringify(org.PublicRepos))

	rate, _, err := gh.RateLimit()
	if err != nil {
		fmt.Printf("Error fetching rate limit: %#v\n\n", err)
	} else {
		fmt.Println("API Rate Limit: remaining", rate.Remaining, "/", rate.Limit)
		//fmt.Printf("API Rate Limit: %#v\n\n", rate)
	}
}

// vi:ts=2
