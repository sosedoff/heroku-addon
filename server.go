package addon

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router   *gin.Engine
	manager  Manager
	manifest *Manifest
}

func New(file string, manager Manager) (*Server, error) {
	manifest, err := readManifest(file)
	if err != nil {
		return nil, err
	}

	server := &Server{
		manifest: manifest,
		manager:  manager,
	}
	server.configure()

	return server, nil
}

func (s *Server) configure() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	group := router.Group("/heroku")
	group.Use(s.basicAuth)
	group.POST("/resources", s.provisionResource)
	group.PUT("/resources/:id", s.modifyResource)
	group.DELETE("/resources/:id", s.deleteResource)

	s.router = router
}

func (s *Server) basicAuth(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.Header("WWW-Authenticate", `Basic realm="Heroku"`)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	auth, err := parseBasicAuth(authHeader)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !s.manifest.isValidAuth(auth) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("auth", auth)
}

func (s *Server) provisionResource(c *gin.Context) {
	req := ProvisionRequest{}

	if err := c.BindJSON(&req); err != nil {
		log.Println("provision request error:", err)
		errorResponse(c, "Unable to parse provision request", 400)
		return
	}

	resource, err := s.manager.Provision(&req)
	if err != nil {
		log.Println("provision request error:", err)
		errorResponse(c, "Unable to provision resource", 400)
		return
	}
	if resource == nil {
		log.Println("provision returned nil resource")
		errorResponse(c, "Unable to provision resource", 400)
		return
	}

	successResponse(c, ProvisionResponse{
		Id:     resource.Id,
		Config: resource.Config,
	})
}

func (s *Server) modifyResource(c *gin.Context) {
	req := &ModifyRequest{}

	if err := c.Bind(req); err != nil {
		errorResponse(c, "Unable to parse modify request", 400)
		return
	}
	req.UUID = c.Param("id")

	resource, err := s.manager.Modify(req)
	if err != nil {
		errorResponse(c, "Unable to modify the resource", 400)
		return
	}
	if resource == nil {
		errorResponse(c, "Resource was not found", 400)
		return
	}

	successResponse(c, ModifyResponse{
		Config:  resource.Config,
		Message: "Resource has been modified",
	})
}

func (s *Server) deleteResource(c *gin.Context) {
	req := &DeleteRequest{
		UUID: c.Param("id"),
	}

	resource, err := s.manager.Delete(req)
	if err != nil {
		errorResponse(c, "Unable to delete the resource", 400)
		return
	}

	successResponse(c, DeleteResponse{
		Message: fmt.Sprintf("Resource %s has been deleted", resource.Id),
	})
}

func (s *Server) Start(bind string) error {
	return s.router.Run(bind)
}
