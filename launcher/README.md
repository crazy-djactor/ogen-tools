# Olympus automatic testnet launcher

# How to use

Requirements:

- Debian based Linux distribution.
- `make` (to execute makefile comands).
- Golang

### Get the program

To get the Compiler you can either build yourself or download it directly from GitHub <https://github.com/olympus-protocol/ogen-tools/releases/latest>

To build it simply use the common golang build command `go build main.go`.

### Flags

| Flag        | Type   | Description                                                                |
|-------------|--------|----------------------------------------------------------------------------|
| `--password`    | string | Password for keystore and wallet                              |
| `--nodes`    | int | Setup the amount of nodes the testnet (minimum of 3 nodes)                              |
| `--validators`    | int | Define the amount of validators per node                            |