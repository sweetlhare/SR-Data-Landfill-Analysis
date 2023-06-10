package httpapi

import (
	"net/http"
	"strconv"
	logicEntities "svalka-service/internal/logic/entities"

	"github.com/gin-gonic/gin"
)

// CreateAudit ...
func (api httpApi) CreateAudit(c *gin.Context) {
	var auditRequest logicEntities.AuditCreateRequest
	if err := c.ShouldBindJSON(&auditRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get user_id from context
	auditRequest.UserID = c.GetUint64("user_id")

	resp, err := api.logic.CreateAudit(c, auditRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteAudit ...
func (api httpApi) DeleteAudit(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = api.logic.DeleteAudit(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Audit deleted successfully")

}

// GetAudit ...
func (api httpApi) GetAudit(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := api.logic.GetAudit(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)

}
