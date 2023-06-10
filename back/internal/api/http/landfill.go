package httpapi

import (
	"net/http"
	"strconv"
	logicEntities "svalka-service/internal/logic/entities"

	"github.com/gin-gonic/gin"
)

// CreateLandfill ...
func (api httpApi) CreateLandfill(c *gin.Context) {
	var landfil logicEntities.LandfillCreate
	if err := c.ShouldBindJSON(&landfil); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := api.logic.CreateLandfill(c, landfil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp})

}

// UpdateLandfill ...
func (api httpApi) UpdateLandfill(c *gin.Context) {
	var landfil logicEntities.LandfillUpdate
	if err := c.ShouldBindJSON(&landfil); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := api.logic.UpdateLandfill(c, landfil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Landfil updated successfully"})

}

// GetLandfills ...
func (api httpApi) GetLandfills(c *gin.Context) {
	regionIdString := c.Query("regionID")
	regionID, err := strconv.ParseUint(regionIdString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	landfills, err := api.logic.GetLandfills(c, regionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, landfills)
}

// GetLandfill ...
func (api httpApi) GetLandfill(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	landfill, err := api.logic.GetLandfill(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, landfill)
}

// DeleteLandfill ...
func (api httpApi) DeleteLandfill(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = api.logic.DeleteLandfill(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Landfill deleted successfully")

}
