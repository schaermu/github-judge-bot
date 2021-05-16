[![codecov](https://codecov.io/gh/schaermu/github-judge-bot/branch/main/graph/badge.svg?token=X6R1PQU7GT)](https://codecov.io/gh/schaermu/github-judge-bot) [![CircleCI](https://circleci.com/gh/schaermu/github-judge-bot.svg?style=shield )](https://circleci.com/gh/schaermu/github-judge-bot)
# GitHub Judge Bot - Your OpenSource evaluation buddy
This bot is a port of an [old node.js](https://github.com/schaermu/repolyzer-slackbot) project of mine. The Go code is probably not state-of-the-art, but this project is supposed to help me get to speed.

Use the bot at your own risk.

# Setting it up
## Project
Copy `config.yaml.example` to `config.yaml`.

## Slack
1) Create a new [Slack app](https://api.slack.com/apps).
2) Copy the Signing Secret from the `App Credentials` section to the config (`slack.signing_secret`).
3) Create a new App-Level token with `connections:write` permissions.
4) Copy the App-Level token to the config (`slack.app_token`).
5) Activate Socket Mode on your app.
6) Switch to *OAuth & Permissions*
7) Add the following Bot Token Scopes to your app:
    ```
    app_mentions:read
    chat:write
    chat:write.customize
    ```
8) Switch to *Event Subscriptions*.
9) Subscribe to the following bot events:
    ```
    app_mention
    ```
10) Go back to *OAuth & Permissions* and install the bot in your workspace.
11) Copy the Bot token to the config (`slack.bot_token`)

## GitHub
1) Create a new [Personal Access Token](https://github.com/settings/tokens) with the `repo:public_repo` permission.
2) Copy the token to the config (`github.access_token`) and set up your username in GitHub.

# Running
You can either run the pre-compiled binary or you can use docker (recommended).

## Binary
Make sure your config file is in the same folder as the executable, then simply start the `github-judge-bot` binary.

## Docker
Starting the bot using docker is simple as well:
```
docker run -v "$(pwd)"/config.yaml:/config.yaml:ro ghcr.io/schaermu/github-judge-bot:latest
```