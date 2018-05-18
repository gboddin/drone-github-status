# drone-github-status

[![Build Status](https://hold-on.nobody.run/api/badges/gboddin/drone-github-status/status.svg)](http://hold-on.nobody.run/drone-github-status)

Drone plugin to work on Github commit status.

The plugin supports getting a list of child repos from Github for parent commit status.

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/github-status
docker build --rm -t gboo/github-status .
```

## Usage

```
   --github-token value  Github OAuth token [$GITHUB_TOKEN, $GITHUB_PR_GITHUB_TOKEN, $PLUGIN_GITHUB_TOKEN]
   --github-query value  Github search query [$GITHUB_STATUS_GITHUB_QUERY, $PLUGIN_GITHUB_QUERY]
   --context value       Status context [$GITHUB_STATUS_CONTEXT, $PLUGIN_CONTEXT]
   --state value         Commit state ( one of pending, error, failure, success ) [$GITHUB_STATUS_STATE, $PLUGIN_STATE]
   --commit-id value     Commit SHA to leave a status for [$GITHUB_STATUS_COMMIT_ID, $PLUGIN_COMMIT_ID, $DRONE_COMMIT_SHA]
   --repo-owner value    Repo owner [$GITHUB_STATUS_REPO_OWNER, $PLUGIN_REPO_OWNER, $DRONE_REPO_OWNER]
   --repo-name value     Repo name [$GITHUB_STATUS_REPO_NAME, $PLUGIN_REPO_NAME, $DRONE_REPO_NAME]
   --link value          Status link [$GITHUB_STATUS_LINK, $PLUGIN_LINK, $DRONE_BUILD_LINK]
```

Execute from the working directory:

```sh
./drone-github-status --github-token=aaeae7ae7ae7ae9ae897eaae97ae97a --action=comment --number=5 --repo-owner=Octocat --repo-name=drone-test --message="Hello world"

```

## Drone plugin

### Simple context

```yaml
pipeline:
  # Simple context
  set-commit-status:
    image: gboo/github-status
    secrets: [ github_token ]
    state: success
    context: drone/subtest/functional-test
```

### Multi context

From Drone in parent build :

```yaml
pipeline:
  # Parent build, set commit status for X subrepos :
  set-commit-status:
    image: gboo/github-status
    github_query: "org:drone-plugins topic:drone-plugin"
    secrets: [ github_token ]
    state: pending
    context: drone/subtest/{{ .FullName }}
    
  trigger-downstream:
    image: gboo/github-search-downstream
    github_query: "org:drone-plugins topic:drone-plugin"
    branch: master
    secrets: [ github_token, drone_server, drone_token ]
    params:
      - DRONE_PARENT_REPO_OWNER=${DRONE_REPO_OWNER}
      - DRONE_PARENT_REPO_NAME=${DRONE_REPO_NAME}
      - DRONE_PARENT_COMMIT_HASH=${DRONE_COMMIT_HASH}
```

From Drone in child build :

```yaml
  # Child build, validate on parent's commit status
  set-parent-commit-status:
    image: gboo/github-status
    secrets: [ github_token ]
    state: success
    context: drone/subtest/${DRONE_REPO}
    repo_owner: ${DRONE_PARENT_REPO_OWNER}
    repo_name: ${DRONE_PARENT_REPO_NAME}
    commit_id: ${DRONE_PARENT_COMMIT_HASH}
```
