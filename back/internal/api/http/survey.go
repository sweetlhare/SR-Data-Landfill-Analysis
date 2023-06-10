package httpapi

import (
	"log"
	"net/http"
	"strconv"
	logicEntities "svalka-service/internal/logic/entities"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateSurvey ...
func (api httpApi) CreateSurvey(c *gin.Context) {
	landfillIDString := c.Param("id")
	landfillID, err := strconv.ParseUint(landfillIDString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var surveyRequest logicEntities.SurveyCreateRequest
	surveyRequest.LandfillID = landfillID

	// Get images ...
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]
	for i := range files {
		file := logicEntities.File{
			FileHeader: files[i],
		}
		surveyRequest.RootImages = append(surveyRequest.RootImages, file)
	}

	// Get images date
	dateString := form.Value["date"]
	date, err := time.Parse(time.RFC3339, dateString[0])
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "parse date failed"})
		return
	}
	surveyRequest.Date = date

	// maxSizeString := form.Value["max_size"]
	// maxSize, err := strconv.ParseUint(maxSizeString[0], 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	surveyRequest.MaxSize = 4096
	// get user_id from context
	surveyRequest.UserID = c.GetUint64("user_id")

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
	id, err := strconv.ParseUint(idString, 10, 64)
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

// DeleteSurvey ...
func (api httpApi) DeleteSurvey(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = api.logic.DeleteSurvey(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Survey deleted successfully")

}
