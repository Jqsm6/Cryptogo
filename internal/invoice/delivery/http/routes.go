package http

import (
	"github.com/gin-gonic/gin"

	"Cryptogo/internal/invoice"
)

func MapInvoiceRoutes(invoiceGroup *gin.RouterGroup, h invoice.Handlers) {
	invoiceGroup.POST(createInvoiceEndpoint, h.Create())
	invoiceGroup.GET(checkInvoiceEndpoint, h.Check())
}
