package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olympus-protocol/ogen-deploy/build"
	"github.com/olympus-protocol/ogen-deploy/config"
)

// Server encapsulates a gin server with the building controllers.
type Server struct {
	server     *gin.Engine
	config     config.Config
	controller *build.Controller
}

// Start starts the gin server for listening webhook events.
func (s *Server) Start() error {
	err := s.server.Run(":" + s.config.Port)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) applyRoutes() {
	api := s.server.Group("/")
	{
		api.POST("/", func(c *gin.Context) { apiWrapper(c, s.controller.Handler) })
	}
	s.server.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
		return
	})
}

// NewServer creates a new server instance.
func NewServer(config config.Config) *Server {
	app := gin.Default()
	app.Use(cors.Default())
	s := &Server{
		controller: build.NewController(),
		config:     config,
		server:     app,
	}
	s.applyRoutes()
	return s
}

func apiWrapper(c *gin.Context, method func() (interface{}, error)) {
	res, err := method()
	responseWrapper(res, err, c)
	return
}

func responseWrapper(data interface{}, err error, c *gin.Context) {
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": data})
	return
}