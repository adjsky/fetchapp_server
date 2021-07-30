package middlewares

import (
	"log"
	"strconv"

	"github.com/adjsky/fetchapp_server/pkg/handlers"
	"github.com/gin-gonic/gin"
)

// Logger middleware logs every incoming call to the server
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.URL.Path, c.Request.Method, c.Request.UserAgent(), c.ClientIP())
	}
}

// EnsureParamIsInt middleware checks that a given parameter is an integer
func EnsureParamIsInt(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := strconv.Atoi(c.Param(param))
		if err != nil {
			handlers.NotFound(c)
			c.Abort()
		}
	}
}
