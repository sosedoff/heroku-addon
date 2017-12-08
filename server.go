package addon

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"time"

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
	group.POST("/sso", s.handleSSO)

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

func (s *Server) handleSSO(c *gin.Context) {
	timestamp := c.Request.FormValue("timestamp")

	preToken := fmt.Sprintf(
		"%s:%s:%s",
		c.Request.FormValue("id"),
		s.manifest.Api.SsoSalt,
		timestamp,
	)
	token := fmt.Sprintf("%x", sha1.Sum([]byte(preToken)))

	if token != c.Request.FormValue("token") {
		c.AbortWithStatus(403)
		return
	}

	var timestampVal int64
	fmt.Sscanf(timestamp, "%d", &timestampVal)

	if timestampVal < time.Now().Unix()-60*2 {
		c.AbortWithStatus(403)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Path:  "/",
		Name:  "heroku-nav-data",
		Value: c.Request.FormValue("nav-data"),
	})

	c.Redirect(302, "/")
}

func (s *Server) Start(bind string) error {
	return s.router.Run(bind)
}
