package http

import (
	"github.com/gin-gonic/gin"

	"Cryptogo/internal/invoice"
)

func MapInvoiceRoutes(group *gin.RouterGroup, h invoice.Handlers) {
	group.POST(createInvoiceEndpoint, h.Create())
	group.GET(checkInvoiceEndpoint, h.Check())
}
