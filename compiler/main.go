package main

import (
	"flag"
	"log"

	"github.com/olympus-protocol/ogen-tools/compiler/config"
	"github.com/olympus-protocol/ogen-tools/compiler/server"
)

func main() {
	// Load the flags
	var datadir, branch, port string
	flag.StringVar(&datadir, "datadir", "/root/ogen-deploy", "Full path of the folder to store the files (will be created if not found).")
	flag.StringVar(&port, "port", "8080", "Define the port for the API request listener.")
	flag.StringVar(&branch, "branch", "master", "Define the branch used to monitor commits and updates.")
	flag.Parse()
	c := config.Config{
		DataDir: datadir,
		Port:    port,
		Branch:  branch,
	}
	s, err := server.NewServer(c)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Start()
	panic(err)
}
