package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
		asModel := sh.statusUC.GetAPIStatus()

		ctx.JSON(http.StatusOK, asModel)
	}
}
