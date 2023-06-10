package httpapi

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"

	"svalka-service/internal/httpserver"
	logicInterfaces "svalka-service/internal/logic/interfaces"
)

type httpApi struct {
	logic logicInterfaces.Logic
	auth  authMiddleware
}

func NewHttpApi(_ context.Context, logic logicInterfaces.Logic) (httpserver.HttpApi, error) {
	if logic == nil {
		return nil, errors.New("logic not initialized")
	}
	return httpApi{
		logic: logic,
		auth:  NewAuthMiddleware(logic),
	}, nil
}

// AddRoutes
func (api httpApi) AddRoutes(e *gin.Engine) error {
	group := e.Group("api")
	// users
	group.POST("/login", api.Login)
	group.POST("/users", api.auth.Admin, api.CreateUser)

	// regions
	group.POST("/regions", api.auth.Admin, api.CreateRegion)
	group.GET("/regions/all", api.auth.All, api.GetAllRegions)

	// landfills
	group.POST("/landfills", api.auth.Admin, api.CreateLandfill)
	group.PATCH("/landfills", api.auth.Admin, api.UpdateLandfill)
	group.GET("/landfills", api.auth.All, api.GetLandfills)
	group.GET("/landfills/:id", api.auth.All, api.GetLandfill)
	group.DELETE("/landfills/:id", api.auth.All, api.DeleteLandfill)

	// violations
	group.POST("/violations", api.auth.Admin, api.CreateViolation)
	group.GET("/violations/all", api.auth.All, api.GetAllViolations)

	// survey
	group.POST("/landfills/:id/survey", api.auth.All, api.CreateSurvey)
	group.GET("/surveys/:id", api.auth.All, api.GetSurvey)
	group.DELETE("/surveys/:id", api.auth.All, api.DeleteSurvey)

	//audit
	group.POST("/audits", api.auth.All, api.CreateAudit)
	group.DELETE("/audits/:id", api.auth.All, api.DeleteAudit)
	group.GET("/audits/:id", api.auth.All, api.GetAudit)
	return nil
}
