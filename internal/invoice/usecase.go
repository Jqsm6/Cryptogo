package invoice

import (
	"context"

	"Cryptogo/internal/models"
)

type UseCase interface {
	Info(ctx context.Context, paymentInfoRequest *models.PaymentInfoRequest) (*models.PaymentInfoResponse, error)
	Create(ctx context.Context, paymentRequest *models.PaymentRequest) (*models.PaymentResponse, error)
	ConfirmETH(ctx context.Context, paymentConfirmRequest *models.PaymentConfirmRequest) (*models.PaymentConfirmResponse, error)
	CheckID(ctx context.Context, id string) (bool, error)
}
