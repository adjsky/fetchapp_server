package services

import "github.com/gin-gonic/gin"

type Service interface {
	Register(*gin.RouterGroup)
	Close()
}
