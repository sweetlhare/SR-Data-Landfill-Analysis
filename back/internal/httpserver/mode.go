package httpserver

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GinModeFromString ...
func GinModeFromString(s string) string {
	switch strings.ToLower(s) {
	case gin.ReleaseMode:
		return gin.ReleaseMode
	case gin.TestMode:
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
