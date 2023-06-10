package httpapi

import (
	"net/http"
	logicEntities "svalka-service/internal/logic/entities"

	"github.com/gin-gonic/gin"
)

// CreateViolation ...
func (api httpApi) CreateViolation(c *gin.Context) {
	var violation logicEntities.Violation
	if err := c.ShouldBindJSON(&violation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := api.logic.CreateViolation(c, violation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Violation created successfully"})
}

// GetAllViolations ...
func (api httpApi) GetAllViolations(c *gin.Context) {
	violations, err := api.logic.GetAllViolations(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, violations)
}
