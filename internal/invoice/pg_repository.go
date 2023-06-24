package invoice

import (
	"context"

	"Cryptogo/internal/models"
)

type Repository interface {
	Create(ctx context.Context, paymentDB *models.PaymentDB) error
	Info(ctx context.Context, id string) (*models.PaymentInfoResponse, error)
	ChangeStatus(ctx context.Context, id string) error
	CheckID(ctx context.Context, id string) (bool, error)
	CheckHash(ctx context.Context, hash string) (bool, error)
	UpdateHash(ctx context.Context, hash, id string) error
	Get(ctx context.Context, id string) (models.InfoDB, error)
}
