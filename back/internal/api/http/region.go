package httpapi

import (
	"net/http"
	logicEntities "svalka-service/internal/logic/entities"

	"github.com/gin-gonic/gin"
)

// CreateRegion ...
func (api httpApi) CreateRegion(c *gin.Context) {
	var region logicEntities.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := api.logic.CreateRegion(c, region)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Region created successfully"})

}

// GetAllRegions ...
func (api httpApi) GetAllRegions(c *gin.Context) {
	regions, err := api.logic.GetAllRegions(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, regions)
}
