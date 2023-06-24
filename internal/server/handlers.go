package server

import (
	"github.com/gin-gonic/gin"

	invoiceHandlers "Cryptogo/internal/invoice/delivery/http"
	invoiceRepository "Cryptogo/internal/invoice/repository"
	invoiceUC "Cryptogo/internal/invoice/usecase"
)

func (s *Server) MapHandlers(g *gin.Engine) {
	iRepo := invoiceRepository.NewInvoiceRepository(s.db, s.log)

	iUC := invoiceUC.NewInvoiceUseCase(iRepo, s.log, s.cfg)

	iHandlers := invoiceHandlers.NewInvoiceHandlers(iUC, s.log)

	versionGroup := g.Group("/v1")
	versionGroup.Use(gin.Logger())

	invoiceHandlers.MapInvoiceRoutes(versionGroup, iHandlers)
}
