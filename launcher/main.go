package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"github.com/olympus-protocol/ogen-tools/launcher/config"
	"github.com/sethvargo/go-password/password"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

var datadir = "./data/"

var ogenSubFolderPrefix = "ogen-node-"

func main() {

	log.Println("Loading Configuration")
	c := loadConfig()

	log.Println("Creating Folder Structure")
	err := folders(c)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Downloading Ogen")
	err = downloadOgen()
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig() config.Config {
	var pass string
	var nodes, validators int

	flag.StringVar(&pass, "password", "", "Password for keystore and wallet")
	flag.IntVar(&nodes, "nodes", 5, "Setup the amount of nodes the testnet (minimum of 5 nodes)")
	flag.IntVar(&validators, "validators", 32, "Define the amount of validators per node (default 32 nodes)")
	flag.Parse()

	if pass == "" {
		pass, _ = password.Generate(32, 10, 0, false, false)
	}

	c := config.Config{
		Password: pass,
		Nodes: nodes,
		Validators: validators,
	}
	return c
}

func folders(c config.Config) error {

	_ = os.RemoveAll(datadir)

	err := os.Mkdir(datadir, 0777)
	if err != nil {
		return err
	}

	for i := 1 ; i <= c.Nodes; i++ {
		numStr := strconv.Itoa(i)
		err := os.Mkdir(path.Join(datadir, ogenSubFolderPrefix + numStr), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadOgen() error {
	_ = os.RemoveAll("./bin")

	file := "https://public.oly.tech/olympus/ogen-release/ogen-0.0.1-linux-amd64.tar.gz"
	resp, err := http.Get(file)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_ = os.Mkdir("./bin", 0777)

	err = extractTar(resp.Body)
	if err != nil {
		return err
	}

	err = os.Rename("./ogen-0.0.1/ogen", "./bin/ogen")
	if err != nil {
		return err
	}

	err = os.Remove("./ogen-0.0.1")
	if err != nil {
		return err
	}

	return nil
}

func extractTar(stream io.Reader) error {
	log.Println("Extracting Ogen")

	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {

		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return err
			}

		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			err = outFile.Close()
			if err != nil {
				return err
			}

		default:
			return err
		}
	}
	return nil
}
