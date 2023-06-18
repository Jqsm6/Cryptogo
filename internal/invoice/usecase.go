package invoice

import (
	"context"

	"Cryptogo/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, prqm *models.PaymentRequest) (*models.PaymentResponse, error)
	Check(ctx context.Context, pirq *models.PaymentInfoRequest) (*models.PaymentInfoResponse, error)
}
