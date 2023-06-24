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
		var paymentRequest *models.PaymentRequest
		var errorResponse models.ErrorResponse

		if err := ctx.BindJSON(&paymentRequest); err != nil {
			errorResponse.Error.Code = http.StatusBadRequest
			errorResponse.Error.Message = "Invalid request body. Use the documentation."
			ctx.JSON(http.StatusBadRequest, &errorResponse)
			return
		}
		if paymentRequest.Currency != "ETH" {
			errorResponse.Error.Code = http.StatusBadRequest
			errorResponse.Error.Message = "At the moment, only 'ETH' is available."
			ctx.JSON(http.StatusBadRequest, &errorResponse)
			return
		}

		resp, err := ch.invoiceUC.Create(ctx, paymentRequest)
		if err != nil {
			errorResponse.Error.Code = http.StatusInternalServerError
			errorResponse.Error.Message = "An error occurred on the server. Retry the request or wait."
			ctx.JSON(http.StatusInternalServerError, &errorResponse)
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

func (ch *invoiceHandlers) Info() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var paymentInfoRequest *models.PaymentInfoRequest
		var errorResponse models.ErrorResponse

		if err := ctx.BindJSON(&paymentInfoRequest); err != nil {
			errorResponse.Error.Code = http.StatusBadRequest
			errorResponse.Error.Message = "Invalid request body. Use the documentation."
			ctx.JSON(http.StatusBadRequest, &errorResponse)
			return
		}

		result, err := ch.invoiceUC.CheckID(ctx, paymentInfoRequest.ID)
		if err != nil {
			errorResponse.Error.Code = http.StatusInternalServerError
			errorResponse.Error.Message = "It was not possible to check the ID for existence. Server side error."
			ctx.JSON(http.StatusInternalServerError, &errorResponse)
			return
		}

		if !result {
			errorResponse.Error.Code = http.StatusBadRequest
			errorResponse.Error.Message = "Invoice with this ID was not found."
			ctx.JSON(http.StatusBadRequest, &errorResponse)
			return
		}

		resp, err := ch.invoiceUC.Info(ctx, paymentInfoRequest)
		if err != nil {
			errorResponse.Error.Code = http.StatusInternalServerError
			errorResponse.Error.Message = "An error occurred on the server. Retry the request or wait."
			ctx.JSON(http.StatusInternalServerError, &errorResponse)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (ch *invoiceHandlers) Confirm() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var paymentConfirmRequest *models.PaymentConfirmRequest
		var errorResponse models.ErrorResponse

		if err := ctx.BindJSON(&paymentConfirmRequest); err != nil {
			errorResponse.Error.Code = http.StatusBadRequest
			errorResponse.Error.Message = "Invalid request body. Use the documentation."
			ctx.JSON(http.StatusBadRequest, &errorResponse)
			return
		}

		paymentConfirmResponse, err := ch.invoiceUC.ConfirmETH(ctx, paymentConfirmRequest)
		if err != nil {
			if err.Error() == "was earlier" {
				errorResponse.Error.Code = http.StatusBadRequest
				errorResponse.Error.Message = "This hash has already been paid."
				ctx.JSON(http.StatusBadRequest, &errorResponse)
				return
			}
			errorResponse.Error.Code = http.StatusInternalServerError
			errorResponse.Error.Message = "An error occurred on the server. Retry the request or wait."
			ctx.JSON(http.StatusInternalServerError, &errorResponse)
			return
		}

		ctx.JSON(http.StatusOK, paymentConfirmResponse)
	}
}
