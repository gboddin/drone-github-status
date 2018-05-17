package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "github status plugin"
	app.Usage = "github status plugin"
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "github-token",
			Usage:  "Github OAuth token",
			EnvVar: "GITHUB_TOKEN,GITHUB_PR_GITHUB_TOKEN,PLUGIN_GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:   "github-query",
			Usage:  "Github search query",
			EnvVar: "GITHUB_STATUS_GITHUB_QUERY,PLUGIN_GITHUB_QUERY",
		},
		cli.StringFlag{
			Name:   "context",
			Usage:  "Status context",
			EnvVar: "GITHUB_STATUS_CONTEXT,PLUGIN_CONTEXT",
		},
		cli.StringFlag{
			Name:   "state",
			Usage:  "Commit state ( one of pending, error, failure, success )",
			EnvVar: "GITHUB_STATUS_STATE,PLUGIN_STATE",
		},
		cli.StringFlag{
			Name:   "commit-id",
			Usage:  "Commit SHA to leave a status for",
			EnvVar: "GITHUB_STATUS_COMMIT_ID,PLUGIN_COMMIT_ID,DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "repo-owner",
			Usage:  "Repo owner",
			EnvVar: "GITHUB_STATUS_REPO_OWNER,PLUGIN_REPO_OWNER,DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo-name",
			Usage:  "Repo name",
			EnvVar: "GITHUB_STATUS_REPO_NAME,PLUGIN_REPO_NAME,DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "link",
			Usage:  "Status link",
			EnvVar: "GITHUB_STATUS_LINK,PLUGIN_LINK,DRONE_BUILD_LINK",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		GithubToken: c.String("github-token"),
		GithubQuery: c.String("github-query"),
		Context:     c.String("context"),
		CommitId:    c.String("commit-id"),
		State:       c.String("state"),
		Link:        c.String("link"),
		RepoName:    c.String("repo-name"),
		RepoOwner:   c.String("repo-owner"),
	}
	return plugin.Exec()
}
