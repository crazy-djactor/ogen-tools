package build

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/olympus-protocol/ogen-tools/compiler/config"
	"github.com/olympus-protocol/ogen-tools/compiler/models"
)

// Controller is the handler for GitHub WebHooks.
type Controller struct {
	config    config.Config
	buildLock sync.RWMutex
}

// move will adjust the produced binaries to match the data directory structure
func (c *Controller) move() error {
	_ = os.RemoveAll(path.Join(c.config.DataDir, "ogen-release"))
	_ = os.Rename("./ogen/release", path.Join(c.config.DataDir, "ogen-release"))
	_ = os.Chmod(path.Join(c.config.DataDir, "ogen-release"), 0777)
	return nil
}

// build will produce new binaries
func (c *Controller) build() error {
	cmd := exec.Command("scripts/build-docker.sh")
	p, err := filepath.Abs("ogen/")
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = p
	err = cmd.Run()
	if err != nil {
		return err
	}
	// Once finished, the binaries are at build/ogen/
	return nil
}

// folder will create the folder path and remove if necesary
func (c *Controller) folder() error {
	_ = os.MkdirAll(c.config.DataDir, 0777)
	_ = os.Remove(path.Join(c.config.DataDir, "ogen-release/"))
	return nil
}

// clone will clone the ogen repository, if it already exists will
func (c *Controller) clone() error {
clone:
	_, err := git.PlainClone("./ogen", false, &git.CloneOptions{
		URL:           "https://github.com/olympus-protocol/ogen",
		Progress:      os.Stdout,
		ReferenceName: plumbing.NewBranchReferenceName(c.config.Branch),
	})
	if err != nil {
		// If the repo already exists, delete it and clone again
		if err == git.ErrRepositoryAlreadyExists {
			_ = os.RemoveAll("./ogen")
			goto clone
		}
		return err
	}
	return nil
}

// Handler is the entrance of the GitHub Webhooks POST information.
func (c *Controller) Handler(payload []byte) (interface{}, error) {
	log.Println("received new push")
	var data models.PushEvent
	err := json.Unmarshal(payload, &data)
	if err != nil {
		log.Println("unable to unmarshal payload data")
		return nil, err
	}
	// If the branch doesn't match "observed" branch, return without issues.
	branch := strings.Split(data.Ref, "/")[2]
	if branch != c.config.Branch {
		log.Println("branch is not observed, skiping")
		return nil, nil
	}
	go func() {
		c.buildLock.Lock()
		defer c.buildLock.Unlock()
		log.Println("start build")
		log.Println("clone repository")
		err = c.clone()
		if err != nil {
			log.Println("error: " + err.Error())
			return
		}
		log.Println("creating folder structure")
		err = c.folder()
		if err != nil {
			log.Println("error: " + err.Error())
			return
		}
		log.Println("building ogen")
		err = c.build()
		if err != nil {
			log.Println("error: " + err.Error())
			return
		}
		log.Println("moving files")
		err = c.move()
		if err != nil {
			log.Println("error: " + err.Error())
			return
		}
	}()
	return nil, nil
}

// NewController returns a new build controller.
func NewController(config config.Config) (*Controller, error) {
	return &Controller{config: config}, nil
}
