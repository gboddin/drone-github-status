package main

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"html/template"
	"math"
	"time"
)

// Plugin defines the Downstream plugin parameters.
type Plugin struct {
	GithubToken string
	GithubQuery string
	Context     string
	State       string
	CommitId    string
	Link        string
	RepoOwner   string
	RepoName    string
}

// Exec runs the plugin
func (p *Plugin) Exec() error {
	if len(p.GithubToken) == 0 {
		return fmt.Errorf("Error: you must provide your Github access token.")
	}

	// Instantiate Github client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.GithubToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	if len(p.GithubQuery) > 0 {
		query := fmt.Sprintf(p.GithubQuery)

		page := 1
		maxPage := math.MaxInt32

		opts := &github.SearchOptions{
			Sort:  "updated",
			Order: "desc",
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}

		for ; page <= maxPage; page++ {
			opts.Page = page
			result, response, err := client.Search.Repositories(oauth2.NoContext, query, opts)
			wait(response)
			if err != nil {
				return err
			}
			maxPage = response.LastPage
			for _, repo := range result.Repositories {
				err := setCommitStatus(
					client,
					p.RepoOwner,
					p.RepoName,
					p.CommitId,
					parseContext(&repo, p.Context),
					p.State,
					p.Link)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	// We're not in multi mode, defaulting to single repo
	repo, _ , err := client.Repositories.Get(oauth2.NoContext, p.RepoOwner, p.RepoName)
	if err != nil {
		return err
	}
	return setCommitStatus(client, p.RepoOwner, p.RepoName, p.CommitId, parseContext(repo, p.Context), p.State, p.Link)
}

func wait(response *github.Response) {
	if response != nil && response.Remaining <= 1 {
		gap := time.Duration(response.Reset.Local().Unix() - time.Now().Unix())
		sleep := gap * time.Second
		if sleep < 0 {
			sleep = -sleep
		}

		time.Sleep(sleep)
	}
}

// Sets a commit status
func setCommitStatus(client *github.Client, repo_owner string, repo_name string, commit_id string, context string, state string, link string) error {
	_, _, err := client.Repositories.CreateStatus(oauth2.NoContext, repo_owner, repo_name, commit_id, &github.RepoStatus{
		Context:   &context,
		State:     &state,
		TargetURL: &link,
	})
	if err == nil {
		fmt.Printf("Context %s set for %s/%s\n", context, repo_owner, repo_name)
	}
	return err
}

func parseContext(repo *github.Repository, contextTemplate string) string {
	// New template from the configured fallback URL
	t, err := template.New("multi-context").Parse(contextTemplate)
	if err != nil {
		logrus.Errorf("Error parsing context template: %s , %s", contextTemplate, err)
	}
	// Prepare return buffer
	context := new(bytes.Buffer)
	// Execute the template with all 3 objects referenced
	t.Execute(context, repo)

	// Return result
	return context.String()
}
