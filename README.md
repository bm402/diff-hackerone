# diff-hackerone
Monitors changes in the HackerOne program directory.

## Features
- Locally stores all programs and assets in the HackerOne public directory
- Queries HackerOne to see whether any programs or assets have been added, updated or deleted
- If changes are found, sends a Slack notification with details of the change and updates the local database
- Designed to be run periodically to get notified of updates as they are released

## Requirements
- Go
- MongoDB Go driver
- MongoDB instance running on `localhost:27017`
- Slack Incoming Webhook URL set as `SLACK_WEBHOOK_URL` environment variable

## Usage
```
$ go run *.go
```

## TODO
- Tidy up JSON parsing of HackerOne response string (surely there is a better way than typecasting `map[string]interface{}` 100 times?)
- Dockerisation to avoid local database dependency
- Depending on the update to the program or asset, give recommendations on which actions to take (eg. for a new URL asset with a wildcard at the start, suggest subdomain enumeration commands). These commands could even be run automatically and the results sent as Slack notifications
