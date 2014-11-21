package main

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/google/go-github/github"
)

func main() {
	//botname := "arlib0"

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: "e282d0ab13c38a4303e65620aeab13c2beba3385"},
	}

	client := github.NewClient(t.Client())

	org, _, err := client.Organizations.Get("arlibs")
	if err != nil {
		fmt.Printf("error: %v\n\n", err)
	} else {
		//fmt.Printf("%v\n\n", github.Stringify(org))
	}

	fmt.Println("Organization: ", *org.Login)
	fmt.Println("Repositories: ", github.Stringify(org.PublicRepos))

	teams, _, err := client.Organizations.ListTeams("arlibs", nil)
	team := teams[0]
	fmt.Println("Teams : ", *team.Name)

	repos, _, err := client.Organizations.ListTeamRepos(*team.ID, nil)
	for _, repo := range repos {
		fmt.Println("Repos : ", *repo.Name, *repo.Description, "URL:", *repo.CloneURL)
	}

	rate, _, err := client.RateLimit()
	if err != nil {
		fmt.Printf("Error fetching rate limit: %#v\n\n", err)
	} else {
		fmt.Println("API Rate Limit: remaining", rate.Remaining, "/", rate.Limit)
		//fmt.Printf("API Rate Limit: %#v\n\n", rate)
	}
}

// vi:ts=4
