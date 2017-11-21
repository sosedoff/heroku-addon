package addon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var exampleManifest = `
{
  "id": "example",
  "api": {
    "config_vars": [
      "EXAMPLE_USER",
      "EXAMPLE_PASSWORD"
    ],
    "regions": [
      "us-west"
    ],
    "password": "password",
    "sso_salt": "salt"
  }
}`

var badProvisionRequest = `
{
  "uuid": "bad"
}`

func makeTestManifest() *Manifest {
	gin.SetMode(gin.TestMode)

	manifest := &Manifest{}
	json.Unmarshal([]byte(exampleManifest), manifest)
	return manifest
}

func TestServer(t *testing.T) {
	manifest := makeTestManifest()
	server := Server{
		manifest: manifest,
		manager:  &TestManager{},
	}
	server.configure()

	// Test empty auth
	req, _ := http.NewRequest("POST", "/heroku/resources", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, 401, resp.Code)
	assert.Equal(t, "", resp.Body.String())
	assert.Equal(t, `Basic realm="Heroku"`, resp.Header().Get("www-authenticate"))

	// Test invalid auth
	req, _ = http.NewRequest("POST", "/heroku/resources", nil)
	resp = httptest.NewRecorder()
	req.SetBasicAuth("foo", "bar")
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, 401, resp.Code)
	assert.Equal(t, "", resp.Body.String())

	// Test valid auth
	req, _ = http.NewRequest("POST", "/heroku/resources", bytes.NewReader(nil))
	resp = httptest.NewRecorder()
	req.SetBasicAuth(manifest.Id, manifest.Api.Password)
	server.router.ServeHTTP(resp, req)
	assert.NotEqual(t, 401, resp.Code)

	// Test bad provision
	req, _ = http.NewRequest("POST", "/heroku/resources", bytes.NewReader([]byte(badProvisionRequest)))
	resp = httptest.NewRecorder()
	req.SetBasicAuth(manifest.Id, manifest.Api.Password)
	server.router.ServeHTTP(resp, req)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, `{"error":"Unable to provision resource","status":400}`, resp.Body.String())
}
