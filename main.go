package main

import (
	"arduino.cc/repository/libraries/db"
	"code.google.com/p/goauth2/oauth"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"strconv"
	s "strings"
	//"os/exec"
	"os"
)

// Configuration part
// TODO: Move into his own file
var gh_auth_token = "e282d0ab13c38a4303e65620aeab13c2beba3385"
var org = "arlibs"
var gh_callbackUrl = "http://kungfu.bug.st:8088/github/event/"
var localGitFolder = "git/"
var librariesIndexFile = os.Getenv("HOME") + "/.arduino15/library_index.json"

// Global github client
var gh_auth = &oauth.Transport{Token: &oauth.Token{AccessToken: gh_auth_token}}
var gh = github.NewClient(gh_auth.Client())

// Make a library by reading library.properties from a github.PullRequest
func MakeLibraryFromPullRequest(pull *github.PullRequest) (*db.LibraryMetadata, error) {
	head := *pull.Head
	headRepo := *head.Repo

	// Get library.properties from pull request HEAD
	getContentOpts := &github.RepositoryContentGetOptions{
		Ref: *head.SHA,
	}
	libPropContent, _, _, err := gh.Repositories.GetContents(*headRepo.Owner.Login, *headRepo.Name, "library.properties", getContentOpts)
	if err != nil || libPropContent == nil {
		fmt.Println("cannot fetch library.properties:" + github.Stringify(err))
		return nil, err
	}
	return MakeLibraryFromRepositoryContent(libPropContent)
}

// Make a library by reading library.properties from a github.RepositoryContent
func MakeLibraryFromRepositoryContent(content *github.RepositoryContent) (*db.LibraryMetadata, error) {
	libPropertiesData, err := base64.StdEncoding.DecodeString(*content.Content)
	if err != nil {
		return nil, err
	}
	return db.ParseLibraryProperties(libPropertiesData)
}

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
		c.String(200, "Pong: "+*ping.Zen)
		return

	case "pull_request":
		var event github.PullRequestEvent
		c.Bind(&event)
		switch *event.Action {
		case "opened", "synchronize":
			go ProcessOpenPullRequest(event.PullRequest)
		case "closed":
			go ProcessClosePullRequest(event.PullRequest)
		}
		c.String(200, "Successfully processed pull_request")
		return

	case "issue_comment":
		var event github.IssueCommentEvent
		c.Bind(&event)
		if event.Issue.PullRequestLinks == nil {
			go ProcessIssueComment(&event)
		} else {
			go ProcessPullRequestComment(&event)
		}
		c.String(200, "Successfully processed issue_comment")
		return
	}

	c.String(200, "Received "+eventType+" from github. Ignoring...")
}

func ProcessIssueComment(event *github.IssueCommentEvent) {
	fmt.Println("Issue comment received")
}

func ProcessPullRequestComment(event *github.IssueCommentEvent) {
	repository := event.Repo
	issue := event.Issue
	pull, _, err := GetPullRequestFromIssue(repository, issue)
	if err != nil {
		fmt.Println("cannot fetch pull request data:" + github.Stringify(err))
		return
	}
	comment := event.Comment
	user := *comment.User.Login
	body := *comment.Body

	// Is an admin?
	// TODO: Check if user is on 'owners' team
	isAdmin := (user == "arlib0")

	if isAdmin {
		body = s.TrimSpace(body)
		switch body {
		case "ok to merge":
			result, _, err := MergePullRequest(pull, user+" accepted "+*pull.Title)
			if err != nil {
				fmt.Println("Error during merge: " + github.Stringify(err))
				return
			}
			fmt.Println(github.Stringify(result))
			return
		}

	}

	fmt.Println(user + " wrote " + body + " on " + *pull.Title)
}

// Get the pull request associated with the issue
func GetPullRequestFromIssue(repository *github.Repository, issue *github.Issue) (*github.PullRequest, *github.Response, error) {
	return gh.PullRequests.Get(*repository.Owner.Login, *repository.Name, *issue.Number)
}

// Merge the pull request
func MergePullRequest(pull *github.PullRequest, commitMessage string) (*github.PullRequestMergeResult, *github.Response, error) {
	repo := pull.Base.Repo
	fmt.Println(*repo.Owner.Login + " / " + *repo.Name + ", " + strconv.Itoa(*pull.Number) + " - " + commitMessage)
	return gh.PullRequests.Merge(*repo.Owner.Login, *repo.Name, *pull.Number, commitMessage)
}

func ProcessOpenPullRequest(pull *github.PullRequest) {
	//commits := *pull.Commits
	title := *pull.Title

	// Is a release request?
	if s.HasPrefix(title, "[RELEASE] ") {

		head := *pull.Head
		base := *pull.Base
		headRepo := *head.Repo
		baseRepo := *base.Repo

		// Check if the release number is the same inside library
		fmt.Println("Doing release!")
		fmt.Println("  Pull request title: '" + title + "'")
		fmt.Println("  release sha       : " + *head.SHA + " in " + *headRepo.FullName)

		// Get library.properties from pull request HEAD and decode library content
		library, err := MakeLibraryFromPullRequest(pull)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		// TODO: Check if the pull request is against master
		// TODO: Check if the pull request is made of only one commit

		// Processing output
		resultMsg := "Hi @" + *pull.User.Login + ",\n"
		resultMsg += "thanks for your submission!\n"
		resultMsg += "\n"
		resultMsg += "Checking library.properties contents for " + *library.Name + "\n"
		errors := 0

		// Check if library name is the same as repository name
		if *library.Name != *baseRepo.Name {
			resultMsg += "  * **ERROR** library 'name' must be " + *baseRepo.Name + " instead of " + *library.Name + "\n"
			errors++
		}
		// Check if pull declared version match the version on manifest file
		version := title[10:]
		if *library.Version != version {
			resultMsg += "  * **ERROR** library 'version' must be " + version + " instead of " + *library.Version + "\n"
			errors++
		}
		// Check author and mainteiner existence
		if library.Author == nil || library.Maintainer == nil {
			resultMsg += "  * **ERROR** 'author' and 'maintainer' fields must be defined\n"
			errors++
		}
		// Check sentence and paragraph and url existence
		if library.Sentence == nil || library.Paragraph == nil || library.URL == nil {
			resultMsg += "  * **ERROR** 'sentence', 'paragraph' and 'url' fields must be defined\n"
			errors++
		}
		// Check architectures
		architectures := s.Split(*library.Architectures, ",")
		for _, arch := range architectures {
			arch = s.TrimSpace(arch)
			resultMsg += "  * Found valid architecture '" + arch + "'\n"
		}

		if errors == 0 {
			resultMsg += "\n"
			resultMsg += "No errors found.\n"
			resultMsg += "\n"
			resultMsg += "This pull request is ready to be merged.\n"
		} else {
			resultMsg += "\n"
			resultMsg += strconv.Itoa(errors) + " errors found.\n"
			resultMsg += "\n"
			resultMsg += "Please fix it and resubmit or update the pullrequest.\n"
		}

		// Send result of analisys as a pull request message
		_, _, err = CommentOnPullRequest(pull, resultMsg)
		if err != nil {
			fmt.Println(github.Stringify(err))
			return
		}
		fmt.Println(resultMsg)
	}
}

func ProcessClosePullRequest(pull *github.PullRequest) {
	// If pull request has been merged..
	if *pull.Merged {
		owner := *pull.Base.Repo.Owner.Login
		repo := *pull.Base.Repo.Name

		// Get library.properties from pull request HEAD and decode library content
		library, err := MakeLibraryFromPullRequest(pull)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		// Create a release for the merged pull request
		release := &github.RepositoryRelease{
			TagName: github.String("v" + *library.Version),
			// master is the default
			// TargetCommittish : "master"
			Name: github.String("Version " + *library.Version),
			Body: pull.Body,
		}
		newRelease, _, err := gh.Repositories.CreateRelease(owner, repo, release)
		if err != nil {
			fmt.Println("Error creating release: " + github.Stringify(err))
			return
		}
		fmt.Println(github.Stringify(newRelease))
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
	libs, err := db.LoadFromFile(librariesIndexFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Loaded %v libraries.\n", len(libs.Libraries))

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
