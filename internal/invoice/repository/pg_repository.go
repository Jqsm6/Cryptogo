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
	row := ir.db.QueryRowContext(ctx, createInvoice, pdm.ID, pdm.State, pdm.Currency, pdm.Amount,
		pdm.ToAddress, pdm.PrivateKey)
	if row.Err() != nil {
		ir.log.Info().Err(row.Err()).Msg("repository")
		return row.Err()
	}

	return nil
}

func (ir *invoiceRepo) Info(ctx context.Context, id string) (*models.PaymentInfoResponse, error) {
	var pirp models.PaymentInfoResponse

	err := ir.db.QueryRowContext(ctx, infoInvoice, id).Scan(&pirp.State, &pirp.ToAddress, &pirp.Amount, &pirp.Currency)
	if err != nil {
		ir.log.Info().Err(err).Msg("repository")
		return nil, err
	}

	return &pirp, nil
}

func (ir *invoiceRepo) ChangeStatus(ctx context.Context, id string) error {
	row := ir.db.QueryRowContext(ctx, changeInvoiceState, id)
	if row.Err() != nil {
		ir.log.Info().Err(row.Err()).Msg("repository")
		return row.Err()
	}

	return nil
}

func (ir *invoiceRepo) CheckID(ctx context.Context, id string) (bool, error) {
	var cm models.CountID
	err := ir.db.QueryRowContext(ctx, checkID, id).Scan(&cm.Count)
	if err != nil {
		ir.log.Info().Err(err).Msg("repository")
		return false, err
	}

	if cm.Count == 1 {
		return true, nil
	}

	return false, nil
}
