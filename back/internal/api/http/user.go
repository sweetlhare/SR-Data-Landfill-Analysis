package httpapi

import (
	"net/http"
	logicEntities "svalka-service/internal/logic/entities"

	"github.com/gin-gonic/gin"
)

// CreateUser ...
func (api httpApi) CreateUser(c *gin.Context) {
	var user logicEntities.UserCreate
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := api.logic.CreateUser(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login ...
func (api httpApi) Login(c *gin.Context) {
	var login logicEntities.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session, err := api.logic.Login(c, login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}
