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
	"github.com/olympus-protocol/ogen-deploy/config"
	"github.com/olympus-protocol/ogen-deploy/models"
)

// Controller is the handler for GitHub WebHooks.
type Controller struct {
	config    config.Config
	buildLock sync.RWMutex
}

// move will adjust the produced binaries to match the datadir structure
func (c *Controller) move() error {
	_ = os.Rename("./ogen/release", path.Join(c.config.Datadir, "ogen-release"))
	return nil
}

// build will produce new binaries
func (c *Controller) build() error {
	cmd := exec.Command("make", "build_cross_docker")
	path, err := filepath.Abs("ogen/")
	if err != nil {
		return err
	}
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		return err
	}
	// Once finished, the binaries are at build/ogen/
	return nil
}

// folder will create the folder path and remove if necesary
func (c *Controller) folder() error {
	_ = os.MkdirAll(c.config.Datadir, 0777)
	_ = os.Remove(path.Join(c.config.Datadir, "ogen-release/"))
	return nil
}

// clone will clone the ogen repository, if it already exists will
func (c *Controller) clone() error {
	_, err := git.PlainClone("./ogen", false, &git.CloneOptions{
		URL:      "https://github.com/olympus-protocol/ogen",
		Progress: os.Stdout,
	})
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			r, err := git.PlainOpen("./ogen")
			if err != nil {
				return err
			}
			err = r.Fetch(&git.FetchOptions{RemoteName: "origin"})
			if err != nil {
				return err
			}
			w, err := r.Worktree()
			if err != nil {
				return err
			}
			err = w.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(c.config.Branch),
			})
			if err != nil {
				return err
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil {
				return err
			}
			return nil
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
			return
		}
		log.Println("creating folder structure")
		err = c.folder()
		if err != nil {
			return
		}
		log.Println("building ogen")
		err = c.build()
		if err != nil {
			return
		}
		log.Println("moving files")
		err = c.move()
		if err != nil {
			return
		}
	}()
	return nil, nil
}

// NewController returns a new build controller.
func NewController(config config.Config) (*Controller, error) {
	return &Controller{config: config}, nil
}
