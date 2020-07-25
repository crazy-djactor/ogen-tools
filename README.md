# Ogen Automatically Deployment Tool

> A tool for continious building/deploying Ogen using GitHub WebHooks.

## Important Note

This service is only configured to work over Linux. There is no future plans to make it work on any other OS.

## Explanation

This tool will make easier to have a continuous building for a production environment for Olympus.

It uses an API to connect to the GitHub webhooks for triggering builds.

To connect the GitHub webhook make sure your API is open to the web and it has a domain configured, once that's ready, please open an issue with the endpoint to connect to.

## How to use

Requirements:

- Debian based Linux distribution.
- Docker.
- Nginx or Apache.
- SSL Certificates.
- Golang (optional).

### Get the program

To get the Ogen Deployment Tool you can either build yourself or download it directly from GitHub <https://github.com/olympus-protocol/ogen-deploy/releases>

To build it simply use the common golang build command `go build main.go`.

### Flags

| Flag        | Type   | Description                                                                |
|-------------|--------|----------------------------------------------------------------------------|
| `--port`    | string | Define the port for the API request listener.                              |
| `--branch`  | string | Define the branch used to monitor commits and updates.                     |
| `--cross`   | bool   | Set to false to disable cross-compiling on all available platforms.        |
| `--datadir` | string | Full path of the folder to store the files (will be created if not found)  |
