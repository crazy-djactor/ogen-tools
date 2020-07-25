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
func (c *Controller) move(commit string) error {
	// TODO
	if c.config.CrossCompile {
		if c.config.Archive {

		} else {

		}
	} else {
		if c.config.Archive {

		} else {

		}
	}
	return nil
}

// build will produce new binaries
func (c *Controller) build() error {
	var cmd *exec.Cmd
	if c.config.CrossCompile {
		cmd = exec.Command("make", "build_cross_docker")
	} else {
		cmd = exec.Command("make", "build")
	}
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
func (c *Controller) folder(commit string) error {
	_ = os.MkdirAll(c.config.Datadir, 0777)
	if c.config.Archive {
		// If is run on archive mode, create a folder with a commit reference.
		os.Mkdir(path.Join(c.config.Datadir, "ogen-release-", commit[0:8]), 0777)
	} else {
		// If not, then remove the folder and create a new one to remove older builds
		_ = os.Remove(path.Join(c.config.Datadir, "ogen-release/"))
		_ = os.Mkdir(path.Join(c.config.Datadir, "ogen-release/"), 0777)
	}
	// For archive mode, the file path is at config.DataDir + ogen-release-COMMIT/
	// For non-archive mode, the file path is at config.DataDir + ogen-release/
	return nil
}

// clone will clone the ogen repository, if it already exists will
func (c *Controller) clone() (string, error) {
	r, err := git.PlainClone("./ogen", false, &git.CloneOptions{
		URL:      "https://github.com/olympus-protocol/ogen",
		Progress: os.Stdout,
	})
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			r, err := git.PlainOpen("./ogen")
			if err != nil {
				return "", err
			}
			err = r.Fetch(&git.FetchOptions{RemoteName: "origin"})
			if err != nil {
				return "", err
			}
			w, err := r.Worktree()
			if err != nil {
				return "", err
			}
			err = w.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(c.config.Branch),
			})
			if err != nil {
				return "", err
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil {
				return "", err
			}
			ref, err := r.Head()
			if err != nil {
				return "", err
			}
			return ref.Hash().String(), nil
		}
		return "", err
	}
	ref, err := r.Head()
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}

// Handler is the entrance of the GitHub Webhooks POST information.
func (c *Controller) Handler(payload []byte) (interface{}, error) {
	log.Println("received new push")
	c.buildLock.Lock()
	defer c.buildLock.Unlock()
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
	log.Println("start build")
	log.Println("clone repository")
	commit, err := c.clone()
	if err != nil {
		return nil, err
	}
	log.Println("creating folder structure")
	err = c.folder(commit)
	if err != nil {
		return nil, err
	}
	log.Println("building ogen")
	err = c.build()
	if err != nil {
		return nil, err
	}
	log.Println("moving files")
	err = c.move(commit)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// NewController returns a new build controller.
func NewController(config config.Config) (*Controller, error) {
	return &Controller{config: config}, nil
}
