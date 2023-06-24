package http

import (
	"github.com/gin-gonic/gin"

	"Cryptogo/internal/invoice"
)

func MapInvoiceRoutes(group *gin.RouterGroup, h invoice.Handlers) {
	group.GET(infoInvoice, h.Info())
	group.GET(confirmInvoice, h.Confirm())
	group.POST(createInvoice, h.Create())
}
