# Olympus automatic compiler

## How to use

Requirements:

- Debian based Linux distribution.
- `make` (to execute makefile comands).
- Docker.
- Nginx or Apache.
- SSL Certificates.
- Golang (optional).

### Get the program

To get the Ogen Deployment Tool you can either build yourself or download it directly from GitHub <https://github.com/olympus-protocol/ogen-tools/compiler/releases>

To build it simply use the common golang build command `go build main.go`.

### Flags

| Flag        | Type   | Description                                                                |
|-------------|--------|----------------------------------------------------------------------------|
| `--port`    | string | Define the port for the API request listener.                              |
| `--branch`  | string | Define the branch used to monitor commits and updates.                     |
| `--datadir` | string | Full path of the folder to store the files (will be created if not found). |
