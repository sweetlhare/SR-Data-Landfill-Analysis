package httpapi

import (
	"net/http"
	logicentities "svalka-service/internal/logic/entities"
	logicInterfaces "svalka-service/internal/logic/interfaces"

	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	logic logicInterfaces.Logic
}

func NewAuthMiddleware(logic logicInterfaces.Logic) authMiddleware {
	return authMiddleware{
		logic: logic,
	}
}

// auth...
func (a authMiddleware) auth(c *gin.Context, validRoles ...logicentities.UserRole) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Session-ID is empty"})
		c.Abort()
		return
	}

	creds, err := a.logic.Auth(sessionID, validRoles...)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("user_id", creds.UserID)
	c.Set("user_role", creds.Role.ToString())
}

// UserAuth ...
func (a authMiddleware) All(c *gin.Context) {
	a.auth(c, logicentities.UserRoleUser, logicentities.UserRoleAdmin)
}

// UserAuth ...
func (a authMiddleware) Admin(c *gin.Context) {
	a.auth(c, logicentities.UserRoleAdmin)
}

// UserAuth ...
func (a authMiddleware) User(c *gin.Context) {
	a.auth(c, logicentities.UserRoleUser)
}
