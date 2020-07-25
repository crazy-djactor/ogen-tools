package build

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/olympus-protocol/ogen-deploy/config"
	"github.com/olympus-protocol/ogen-deploy/models"
)

// Controller is the handler for GitHub WebHooks.
type Controller struct {
	config config.Config
}

// Handler is the entrance of the GitHub Webhooks POST information.
func (c *Controller) Handler(payload []byte) (interface{}, error) {
	var data models.PushEvent
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return nil, err
	}
	log.Println("received new push")
	// If the branch doesn't match "observed" branch, return without issues.
	branch := strings.Split(data.Ref, "/")[2]
	if branch != c.config.Branch {
		log.Println("branch is not observed, skiping")
		return nil, nil
	}
	log.Println("start build")

	return nil, nil
}

// NewController returns a new build controller.
func NewController(config config.Config) *Controller {
	return &Controller{config: config}
}
