package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Cryptogo/internal/invoice"
	"Cryptogo/internal/models"
	"Cryptogo/pkg/logger"
)

type invoiceHandlers struct {
	invoiceUC invoice.UseCase
	log       *logger.Logger
}

func NewInvoiceHandlers(invoiceUC invoice.UseCase, log *logger.Logger) invoice.Handlers {
	return &invoiceHandlers{invoiceUC: invoiceUC, log: log}
}

func (ch *invoiceHandlers) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var prqm *models.PaymentRequest

		if err := ctx.BindJSON(&prqm); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if prqm.Currency != "ETH" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "only 'ETH' is available"})
			return
		}

		resp, err := ch.invoiceUC.Create(ctx, prqm)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to create invoice"})
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

func (ch *invoiceHandlers) Info() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pirq *models.PaymentInfoRequest

		if err := ctx.BindJSON(&pirq); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		resp, err := ch.invoiceUC.Info(ctx, pirq)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "there was an error, or the invoice could not be found."})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
