package httpserver

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// setupCors ...
func setupCors(r *gin.Engine) {
	r.Use(
		cors.New(cors.Config{
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders: []string{
				"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-Session-ID",
			},
			AllowOriginFunc: func(origin string) bool {
				return true
			},
		}),
	)
}
