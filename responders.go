package addon

import (
	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, message interface{}, status int) {
	c.AbortWithStatusJSON(status, map[string]interface{}{
		"error":  message,
		"status": status,
	})
}

func successResponse(c *gin.Context, data interface{}) {
	c.JSON(200, data)
}
