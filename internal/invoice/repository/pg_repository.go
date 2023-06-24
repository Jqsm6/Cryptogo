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

func (ir *invoiceRepo) Create(ctx context.Context, paymentDB *models.PaymentDB) error {
	row := ir.db.QueryRowContext(ctx, create, paymentDB.ID, paymentDB.State, paymentDB.Currency, paymentDB.Amount,
		paymentDB.ToAddress, paymentDB.FromAddress)
	if row.Err() != nil {
		ir.log.Err(row.Err()).Msg("repository")
		return row.Err()
	}

	return nil
}

func (ir *invoiceRepo) Info(ctx context.Context, id string) (*models.PaymentInfoResponse, error) {
	var paymentInfoResponse models.PaymentInfoResponse

	err := ir.db.QueryRowContext(ctx, info, id).Scan(&paymentInfoResponse.ID, &paymentInfoResponse.State, 
		&paymentInfoResponse.ToAddress, &paymentInfoResponse.Amount, &paymentInfoResponse.Currency, &paymentInfoResponse.FromAddress)
	if err != nil {
		ir.log.Err(err).Msg("repository")
		return nil, err
	}

	return &paymentInfoResponse, nil
}

func (ir *invoiceRepo) ChangeStatus(ctx context.Context, id string) error {
	row := ir.db.QueryRowContext(ctx, changeState, id)
	if row.Err() != nil {
		ir.log.Err(row.Err()).Msg("repository")
		return row.Err()
	}

	return nil
}

func (ir *invoiceRepo) CheckID(ctx context.Context, id string) (bool, error) {
	var count models.Count
	err := ir.db.QueryRowContext(ctx, checkID, id).Scan(&count.Count)
	if err != nil {
		ir.log.Err(err).Msg("repository")
		return false, err
	}

	if count.Count == 1 {
		return true, nil
	}

	return false, nil
}

func (ir *invoiceRepo) CheckHash(ctx context.Context, hash string) (bool, error) {
	var count models.Count
	err := ir.db.QueryRowContext(ctx, checkHash, hash).Scan(&count.Count)
	if err != nil {
		ir.log.Err(err).Msg("repository")
		return false, err
	}

	if count.Count == 1 {
		return true, nil
	}

	return false, nil
}

func (ir *invoiceRepo) UpdateHash(ctx context.Context, hash, id string) error {
	row := ir.db.QueryRowContext(ctx, updateHash, hash, id)
	if row.Err() != nil {
		ir.log.Err(row.Err()).Msg("repository")
		return row.Err()
	}

	return nil
}

func (ir *invoiceRepo) Get(ctx context.Context, id string) (models.InfoDB, error) {
	var info models.InfoDB
	err := ir.db.QueryRowContext(ctx, get, id).Scan(&info.FromAddress, &info.ToAddress, &info.Amount)
	if err != nil {
		ir.log.Err(err).Msg("repository")
		return info, err
	}

	return info, nil
}
