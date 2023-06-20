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
	case "ETH":
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			cuc.log.Info().Err(err).Msg("usecase")
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
		State:       "notpayed",
		Currency:    prqm.Currency,
		Amount:      prqm.Amount,
		ToAddress:   address,
		PrivateKey:  hexutil.Encode(privateKeyBytes),
	}

	err := cuc.repo.Create(ctx, &pdm)
	if err != nil {
		return nil, err
	}

	return &prpm, nil
}

func (cuc *invoiceUseCase) Info(ctx context.Context, pirq *models.PaymentInfoRequest) (*models.PaymentInfoResponse, error) {
	id := pirq.ID

	pirp, err := cuc.repo.Info(ctx, id)
	if err != nil {
		return nil, err
	}

	if pirp.State == "paid" {
		return pirp, nil
	}

	switch pirp.Currency {
	case "ETH":
		result, err := infoETH(pirp)
		if err != nil {
			cuc.log.Info().Err(err).Msg("usecase")
			return nil, err
		}

		if result {
			err := cuc.repo.ChangeStatus(ctx, id)
			if err != nil {
				return nil, err
			}
			pirp.State = "paid"
		}
	}

	return pirp, nil
}

func (cuc *invoiceUseCase) CheckID(ctx context.Context, id string) (bool, error) {
	result, err := cuc.repo.CheckID(ctx, id)
	if err != nil {
		return false, err
	}

	return result, nil
}
