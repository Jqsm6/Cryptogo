package server

import (
	"github.com/gin-gonic/gin"

	invoiceRepository "Cryptogo/internal/invoice/repository"
	invoiceHandlers "Cryptogo/internal/invoice/delivery/http"
	invoiceUC "Cryptogo/internal/invoice/usecase"
	
	statusHandlers "Cryptogo/internal/status/delivery/http"
	statusUC "Cryptogo/internal/status/usecase"
)

func (s *Server) MapHandlers(g *gin.Engine) {
	iRepo := invoiceRepository.NewInvoiceRepository(s.db, s.log)

	iUC := invoiceUC.NewInvoiceUseCase(iRepo, s.log)
	sUC := statusUC.NewStatusUseCase(s.log)

	iHandlers := invoiceHandlers.NewInvoiceHandlers(iUC, s.log)
	sHandlers := statusHandlers.NewStatusHandlers(sUC, s.log)

	invoiceGroup := g.Group("/invoice")
	statusGroup := g.Group("/status")
	statusGroup.Use(gin.Logger())
	invoiceGroup.Use(gin.Logger())

	invoiceHandlers.MapInvoiceRoutes(invoiceGroup, iHandlers)
	statusHandlers.MapStatusRoutes(statusGroup, sHandlers)
}
