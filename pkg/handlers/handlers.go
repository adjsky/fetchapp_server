package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFound is called when a client tries to access a nonexistent resource
func NotFound(c *gin.Context) {
	code := http.StatusNotFound
	c.JSON(code, gin.H{
		"code":    code,
		"message": "no resource found",
	})
}

// NoMethod is called when a client tries to access a resource by not implemented method
func NoMethod(c *gin.Context) {
	code := http.StatusMethodNotAllowed
	c.JSON(code, gin.H{
		"code":    code,
		"message": "method is not allowed",
	})
}
