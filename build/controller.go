package build

// Controller is the handler for GitHub WebHooks.
type Controller struct {
}

// Handler is the entrance of the GitHub Webhooks POST information.
func (c *Controller) Handler() (interface{}, error) {
	return nil, nil
}

// NewController returns a new build controller.
func NewController() *Controller {
	return &Controller{}
}