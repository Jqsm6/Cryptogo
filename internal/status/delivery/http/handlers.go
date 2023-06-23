package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Cryptogo/internal/models"
	"Cryptogo/internal/status"
	"Cryptogo/pkg/logger"
)

type statusHandlers struct {
	statusUC status.UseCase
	log      *logger.Logger
}

func NewStatusHandlers(statusUC status.UseCase, log *logger.Logger) status.Handlers {
	return &statusHandlers{statusUC: statusUC, log: log}
}

func (sh *statusHandlers) GetAPIStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var errorResponse *models.ErrorResponse
		status, err := sh.statusUC.GetAPIStatus()
		if err != nil {
			errorResponse.Error.Code = http.StatusInternalServerError
			errorResponse.Error.Message = "An error occurred on the server. Retry the request or wait."
			ctx.JSON(http.StatusInternalServerError, &errorResponse)
			return
		}

		ctx.JSON(http.StatusOK, status)
	}
}
