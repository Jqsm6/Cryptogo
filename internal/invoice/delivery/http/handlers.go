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
		var em models.ErrorResponse

		if err := ctx.BindJSON(&prqm); err != nil {
			em.Error.Code = http.StatusBadRequest
			em.Error.Message = "Invalid request body. Use the documentation."
			ctx.JSON(http.StatusBadRequest, &em)
			return
		}
		if prqm.Currency != "ETH" {
			em.Error.Code = http.StatusBadRequest
			em.Error.Message = "At the moment, only 'ETH' is available."
			ctx.JSON(http.StatusBadRequest, &em)
			return
		}

		resp, err := ch.invoiceUC.Create(ctx, prqm)
		if err != nil {
			em.Error.Code = http.StatusInternalServerError
			em.Error.Message = "An error occurred on the server. Retry the request or wait."
			ctx.JSON(http.StatusInternalServerError, &em)
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

func (ch *invoiceHandlers) Info() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pirq *models.PaymentInfoRequest
		var em models.ErrorResponse

		if err := ctx.BindJSON(&pirq); err != nil {
			em.Error.Code = http.StatusBadRequest
			em.Error.Message = "Invalid request body. Use the documentation."
			ctx.JSON(http.StatusBadRequest, &em)
			return
		}

		result, err := ch.invoiceUC.CheckID(ctx, pirq.ID)
		if err != nil {
			em.Error.Code = http.StatusInternalServerError
			em.Error.Message = "It was not possible to check the ID for existence. Server side error."
			ctx.JSON(http.StatusInternalServerError, &em)
			return
		}

		if !result {
			em.Error.Code = http.StatusBadRequest
			em.Error.Message = "Invoice with this ID was not found."
			ctx.JSON(http.StatusBadRequest, &em)
			return
		}

		resp, err := ch.invoiceUC.Info(ctx, pirq)
		if err != nil {
			em.Error.Code = http.StatusInternalServerError
			em.Error.Message = "An error occurred on the server. Retry the request or wait."
			ctx.JSON(http.StatusInternalServerError, &em)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
