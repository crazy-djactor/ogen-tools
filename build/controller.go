package build

import (
	"encoding/json"
	"fmt"

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
	fmt.Println(data)
	return nil, nil
}

// NewController returns a new build controller.
func NewController(config config.Config) *Controller {
	return &Controller{config: config}
}
