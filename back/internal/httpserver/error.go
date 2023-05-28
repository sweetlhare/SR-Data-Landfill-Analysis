package httpserver

import (
	"net/http"
	"strings"
	"svalka-service/pkg/custom"

	"github.com/gin-gonic/gin"
)

const (
	FailedAddRoutesError        custom.Error = "failed to add routes"
	FailedStartServerError      custom.Error = "failed to start server"
	FailedListeningError        custom.Error = "failed listening and processing"
	ServerForcedToShutdownError custom.Error = "server forced to shutdown"
	ServerConfigError           custom.Error = "failed to get http server config"
)

// ServerError ....
type ServerError struct {
	Error        string   `json:"error"`
	ErrorDetails []string `json:"error_details"`
}

// SendBindError ...
func SendBindError(ctx *gin.Context, err error) {
	serverError := ServerError{}
	serverError.Error = "Bad request"
	errs := strings.Split(err.Error(), "\n")
	for _, ers := range errs {
		es := strings.Split(ers, "Error:")
		if len(es) > 1 {
			serverError.ErrorDetails = append(serverError.ErrorDetails, es[1])
		}
	}
	ctx.JSON(http.StatusBadRequest, serverError)
}

// SendRespond ...
func SendRespond(ctx *gin.Context, respond interface{}, err error) {
	if err != nil {
		serverError := ServerError{}
		serverError.Error = "Internal server error"
		ctx.JSON(http.StatusInternalServerError, serverError)
	} else {
		ctx.JSON(http.StatusOK, respond)
	}
}

// Bind
func Bind(ctx *gin.Context, any interface{}) error {
	err := ctx.ShouldBind(any)
	if err != nil {
		ctx.Errors = append(ctx.Errors, &gin.Error{Err: err, Type: gin.ErrorTypeBind})
		return err
	}
	return nil
}
