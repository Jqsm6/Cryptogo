package usecase

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/xid"

	"Cryptogo/internal/invoice"
	"Cryptogo/internal/models"
	"Cryptogo/pkg/logger"
)

type invoiceUseCase struct {
	repo invoice.Repository
	log  *logger.Logger
}

func NewInvoiceUseCase(repo invoice.Repository, log *logger.Logger) invoice.UseCase {
	return &invoiceUseCase{repo: repo, log: log}
}

func (cuc *invoiceUseCase) Create(ctx context.Context, prqm *models.PaymentRequest) (*models.PaymentResponse, error) {
	var prpm models.PaymentResponse

	guid := xid.New().String()

	var privateKeyBytes []byte
	var address string

	switch prqm.Currency {
	case "eth":
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			cuc.log.Info().Err(err)
			return nil, err
		}

		privateKeyBytes = crypto.FromECDSA(privateKey)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			cuc.log.Info().Msg("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
			return nil, err
		}
		address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		prpm.ID = guid
		prpm.ToAddress = address
	}

	pdm := models.PaymentDB{
		ID:          guid,
		Status:      0,
		Currency:    prqm.Currency,
		Amount:      prqm.Amount,
		FromAddress: prqm.FromAddress,
		ToAddress:   address,
		PrivateKey:  hexutil.Encode(privateKeyBytes),
	}

	err := cuc.repo.Create(ctx, &pdm)
	if err != nil {
		cuc.log.Info().Err(err)
		return nil, err
	}

	return &prpm, nil
}

func (cuc *invoiceUseCase) Check(ctx context.Context, pirq *models.PaymentInfoRequest) (*models.PaymentInfoResponse, error) {
	id := pirq.ID

	pirp, err := cuc.repo.Check(ctx, id)
	if err != nil {
		return nil, err
	}

	if pirp.Status == 1 {
		return pirp, nil
	}

	switch pirp.Currency {
	case "eth":
		result, err := checkETHReplenishment(pirp)
		if err != nil {
			return nil, err
		}

		if result {
			err := cuc.repo.ChangeStatus(ctx, id)
			if err != nil {
				return nil, err
			}
			pirp.Status = 1
		}
	}

	return pirp, nil
}
