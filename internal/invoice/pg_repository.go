package invoice

import (
	"context"

	"Cryptogo/internal/models"
)

type Repository interface {
	Create(ctx context.Context, pdm *models.PaymentDB) error
	Info(ctx context.Context, id string) (*models.PaymentInfoResponse, error)
	ChangeStatus(ctx context.Context, id string) error
	CheckID(ctx context.Context, id string) (bool, error)
}
