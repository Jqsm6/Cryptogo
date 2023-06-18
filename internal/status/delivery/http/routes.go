package http

import (
	"github.com/gin-gonic/gin"

	"Cryptogo/internal/status"
)

func MapStatusRoutes(statusGroup *gin.RouterGroup, h status.Handlers) {
	statusGroup.GET(statusAPIEndpoint, h.GetAPIStatus())
}
