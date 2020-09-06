# Olympus automatic testnet launcher

# How to use

Requirements:

- Debian based Linux distribution.
- Golang

### Get the program

To get the Compiler you can either build yourself or download it directly from GitHub <https://github.com/olympus-protocol/ogen-tools/releases/latest>

To build the tools run `make build` on the main folder.

### Flags

| Flag        | Type   | Description                                                                |
|-------------|--------|----------------------------------------------------------------------------|
| `--password`    | string | Password for keystore and wallet                              |
| `--nodes`    | int | Setup the amount of nodes the testnet (minimum of 3 nodes)                              |
| `--validators`    | int | Define the amount of validators per node                            |
| `--host`    | string | The external IP address to specify on chain file                           |
| `--source`    | bool | Use this flag to build from source                           |
| `--debug`    | bool | Use this flag to start nodes on debug mode                           |
| `--branch`    | string | When using the `source` you can specify a branch to build from                           |