package http

import (
	"github.com/gin-gonic/gin"

	"Cryptogo/internal/status"
)

func MapStatusRoutes(group *gin.RouterGroup, h status.Handlers) {
	group.GET(statusAPI, h.GetAPIStatus())
}
