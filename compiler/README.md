# Olympus automatic compiler

## How to use

Requirements:

- Debian based Linux distribution.
- Docker.
- Nginx or Apache.
- SSL Certificates.
- Golang (optional).

Once your API is running with a domain with a SSL certificate, please open an issue with the endpoint to connect the GitHub WebHook API to you.

### Get the program

To get the Compiler you can either build yourself or download it directly from GitHub <https://github.com/olympus-protocol/ogen-tools/releases/latest>

To build the tools run `make build` on the main folder.

### Flags

| Flag        | Type   | Description                                                                |
|-------------|--------|----------------------------------------------------------------------------|
| `--port`    | string | Define the port for the API request listener.                              |
| `--branch`  | string | Define the branch used to monitor commits and updates.                     |
| `--datadir` | string | Full path of the folder to store the files (will be created if not found). |
