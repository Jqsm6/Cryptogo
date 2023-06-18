package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"Cryptogo/internal/invoice"
	"Cryptogo/internal/models"
	"Cryptogo/pkg/logger"
)

type invoiceRepo struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewInvoiceRepository(db *sqlx.DB, log *logger.Logger) invoice.Repository {
	return &invoiceRepo{db: db, log: log}
}

func (ir *invoiceRepo) Create(ctx context.Context, pdm *models.PaymentDB) error {
	row := ir.db.QueryRowContext(ctx, createInvoice, pdm.ID, pdm.Status, pdm.Currency, pdm.Amount, pdm.FromAddress,
		pdm.ToAddress, pdm.PrivateKey)
	if row.Err() != nil {
		ir.log.Info().Err(row.Err())
		return row.Err()
	}

	return nil
}

func (ir *invoiceRepo) Check(ctx context.Context, id string) (*models.PaymentInfoResponse, error) {
	var pirp models.PaymentInfoResponse

	err := ir.db.QueryRowContext(ctx, checkInvoice, id).Scan(&pirp.Status, &pirp.ToAddress, &pirp.Amount, &pirp.Currency)
	if err != nil {
		ir.log.Info().Err(err)
		return nil, err
	}

	return &pirp, nil
}

func (ir *invoiceRepo) ChangeStatus(ctx context.Context, id string) error {
	row := ir.db.QueryRowContext(ctx, changeInvoiceStatusEndpoint, id)
	if row.Err() != nil {
		ir.log.Info().Err(row.Err())
		return row.Err()
	}

	return nil
}
