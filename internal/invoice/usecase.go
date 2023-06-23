package invoice

import (
	"context"

	"Cryptogo/internal/models"
)

type UseCase interface {
	Info(ctx context.Context, paymentInfoRequest *models.PaymentInfoRequest) (*models.PaymentInfoResponse, error)
	Create(ctx context.Context, paymentRequest *models.PaymentRequest) (*models.PaymentResponse, error)
	InfoETH(paymentInfoResponse *models.PaymentInfoResponse) (bool, string, error)
	InfoBTC(paymentInfoResponse *models.PaymentInfoResponse) (bool, string, error)
	InfoBNB(paymentInfoResponse *models.PaymentInfoResponse) (bool, string, error)
	CheckID(ctx context.Context, id string) (bool, error)
}
