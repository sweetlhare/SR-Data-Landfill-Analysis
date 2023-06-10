package httpserver

import "github.com/gin-gonic/gin"

type HttpApi interface {
	AddRoutes(e *gin.Engine) error
}
