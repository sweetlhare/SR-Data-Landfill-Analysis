package httpapi

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"svalka-service/internal/httpserver"
	logicEntities "svalka-service/internal/logic/entities"
	logicInterfaces "svalka-service/internal/logic/interfaces"
)

type httpApi struct {
	logic logicInterfaces.Logic
}

func NewHttpApi(_ context.Context, logic logicInterfaces.Logic) (httpserver.HttpApi, error) {
	if logic == nil {
		return nil, errors.New("logic not initialized")
	}
	return httpApi{
		logic: logic,
	}, nil
}

// AddRoutes
func (api httpApi) AddRoutes(e *gin.Engine) error {
	group := e.Group("api")
	// users
	group.POST("/users", api.CreateUser)
	group.POST("/login", api.Login)
	// regions
	group.POST("/regions", api.CreateRegion)
	group.GET("/regions/all", api.GetAllRegions)
	// landfills
	group.POST("/landfills", api.CreateLandfill)
	group.GET("/landfills", api.GetLandfills)
	group.GET("/landfills/:id", api.GetLandfill)
	// violations
	group.POST("/violations", api.CreateViolation)
	group.GET("/violations/all", api.GetAllViolations)
	// survey
	group.POST("/landfills/:id/survey", api.CreateSurvey)
	group.GET("/surveys/:id", api.GetSurvey)
	//audit
	group.POST("/audits", api.CreateAudit)

	return nil
}

// CreateUser ...
func (api httpApi) CreateUser(c *gin.Context) {
	var user logicEntities.User
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

// CreateLandfill ...
func (api httpApi) CreateLandfill(c *gin.Context) {
	var landfil logicEntities.Landfill
	if err := c.ShouldBindJSON(&landfil); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := api.logic.CreateLandfill(c, landfil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Landfil created successfully"})

}

// GetLandfills ...
func (api httpApi) GetLandfills(c *gin.Context) {
	regionIdString := c.Query("regionID")
	regionID, err := strconv.ParseInt(regionIdString, 10, 64)
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

// CreateSurvey ...
func (api httpApi) CreateSurvey(c *gin.Context) {
	landfillIDString := c.Param("id")
	landfillID, err := strconv.ParseInt(landfillIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var surveyRequest logicEntities.SurveyRootRequest
	surveyRequest.LandfillID = landfillID

	// MultipartForm ...
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]
	for i := range files {
		var file logicEntities.File
		file = logicEntities.File{files[i]}
		surveyRequest.RootImages = append(surveyRequest.RootImages, file)
	}

	dateString := form.Value["date"]
	date, err := time.Parse(time.RFC3339, dateString[0])
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "parse date failed"})
		return
	}
	surveyRequest.Date = date

	survey, err := api.logic.CreateSurvey(c, surveyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, survey)
}

// GetSurvey ...
func (api httpApi) GetSurvey(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	survey, err := api.logic.GetSurvey(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, survey)

}

// CreateAudit ...
func (api httpApi) CreateAudit(c *gin.Context) {
	var auditRequest logicEntities.AuditRequest
	if err := c.ShouldBindJSON(&auditRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sessionID := c.GetHeader("X-Session-ID")
	auditRequest.SessionID = sessionID
	survey, err := api.logic.CreateAudit(c, auditRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, survey)
}

// GetLandfill ...
func (api httpApi) GetLandfill(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseInt(idString, 10, 64)
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
