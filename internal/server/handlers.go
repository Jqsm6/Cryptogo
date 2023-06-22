package server

import (
	"github.com/gin-gonic/gin"

	invoiceHandlers "Cryptogo/internal/invoice/delivery/http"
	invoiceRepository "Cryptogo/internal/invoice/repository"
	invoiceUC "Cryptogo/internal/invoice/usecase"

	statusHandlers "Cryptogo/internal/status/delivery/http"
	statusUC "Cryptogo/internal/status/usecase"
)

func (s *Server) MapHandlers(g *gin.Engine) {
	iRepo := invoiceRepository.NewInvoiceRepository(s.db, s.log)

	iUC := invoiceUC.NewInvoiceUseCase(iRepo, s.log, s.cfg)
	sUC := statusUC.NewStatusUseCase(s.log, s.cfg)

	iHandlers := invoiceHandlers.NewInvoiceHandlers(iUC, s.log)
	sHandlers := statusHandlers.NewStatusHandlers(sUC, s.log)

	defaultGroup := g.Group("/")
	versionGroup := g.Group("/v1")
	versionGroup.Use(gin.Logger())
	defaultGroup.Use(gin.Logger())

	invoiceHandlers.MapInvoiceRoutes(versionGroup, iHandlers)
	statusHandlers.MapStatusRoutes(defaultGroup, sHandlers)
}
