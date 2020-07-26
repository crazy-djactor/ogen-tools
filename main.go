package main

import (
	"flag"
	"log"

	"github.com/olympus-protocol/ogen-deploy/config"
	"github.com/olympus-protocol/ogen-deploy/server"
)

func main() {
	// Load the flags
	var datadir, branch, port string
	flag.StringVar(&datadir, "datadir", "/root/ogen-deploy", "Full path of the folder to store the files (will be created if not found).")
	flag.StringVar(&port, "port", "8080", "Define the port for the API request listener.")
	flag.StringVar(&branch, "branch", "master", "Define the branch used to monitor commits and updates.")
	flag.Parse()
	config := config.Config{
		Datadir:      datadir,
		Port:         port,
		Branch:       branch,
	}
	s, err := server.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}
