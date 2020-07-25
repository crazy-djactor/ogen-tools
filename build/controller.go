package build

import (
	"encoding/json"
	"fmt"
)

// Controller is the handler for GitHub WebHooks.
type Controller struct {
}

// Handler is the entrance of the GitHub Webhooks POST information.
func (c *Controller) Handler(payload []byte) (interface{}, error) {
	var data interface{}
	err := json.Unmarshal(payload, data)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	return nil, nil
}

// NewController returns a new build controller.
func NewController() *Controller {
	return &Controller{}
}
