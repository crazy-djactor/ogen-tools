package main

import (
	"flag"

	"github.com/olympus-protocol/ogen-deploy/config"
	"github.com/olympus-protocol/ogen-deploy/server"
)

func main() {
	// Load the flags
	var datadir, branch, port string
	var cross bool
	flag.StringVar(&datadir, "datadir", "/root/ogen-deploy", "Full path of the folder to store the files (will be created if not found)")
	flag.StringVar(&port, "port", "8080", "Define the port for the API request listener.")
	flag.StringVar(&branch, "branch", "master", "Define the branch used to monitor commits and updates.")
	flag.BoolVar(&cross, "cross", true, "Use all to cross-compile or leave it to build for current OS.")
	flag.Parse()
	config := config.Config{
		Datadir:      datadir,
		Port:         port,
		Branch:       branch,
		CrossCompile: cross,
	}
	s := server.NewServer(config)
	s.Start()
}
