package addon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_errorResponse(t *testing.T) {
	router := gin.New()

	router.GET("/error", func(c *gin.Context) {
		errorResponse(c, "Error message", 400)
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, `{"error":"Error message","status":400}`, resp.Body.String())
}

func Test_successResponse(t *testing.T) {
	router := gin.New()
	router.GET("/success", func(c *gin.Context) {
		data := map[string]string{"hello": "world"}
		successResponse(c, data)
	})

	req, _ := http.NewRequest("GET", "/success", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, `{"hello":"world"}`, resp.Body.String())
}
