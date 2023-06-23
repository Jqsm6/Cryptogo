package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/xid"
	"github.com/xorcare/blockchain"

	"Cryptogo/config"
	"Cryptogo/internal/invoice"
	"Cryptogo/internal/models"
	"Cryptogo/pkg/logger"
)

type invoiceUseCase struct {
	repo     invoice.Repository
	log      *logger.Logger
	cfg      *config.Config
	bcClient *blockchain.Client
}

func NewInvoiceUseCase(repo invoice.Repository, log *logger.Logger, cfg *config.Config, bcClient *blockchain.Client) invoice.UseCase {
	return &invoiceUseCase{repo: repo, log: log, cfg: cfg, bcClient: bcClient}
}

func (cuc *invoiceUseCase) Create(ctx context.Context, prqm *models.PaymentRequest) (*models.PaymentResponse, error) {
	var prpm models.PaymentResponse

	guid := xid.New().String()

	switch prqm.Currency {
	case "ETH":
		prpm.ToAddress = cuc.cfg.ETHRecipient
	case "BTC":
		prpm.ToAddress = cuc.cfg.BTCRecipient
	}

	prpm.ID = guid

	strAmount := strconv.FormatFloat(prqm.Amount, 'f', -1, 64)
	pdm := models.PaymentDB{
		ID:          guid,
		State:       "notpayed",
		Currency:    prqm.Currency,
		Amount:      strAmount,
		ToAddress:   prpm.ToAddress,
		FromAddress: prqm.FromAddress,
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
		result, hash, err := cuc.InfoETH(pirp)
		if err != nil {
			return nil, err
		}

		if result {
			err = cuc.repo.UpdateTransactionHash(ctx, hash, pirp.ID)
			if err != nil {
				return nil, err
			}

			err := cuc.repo.ChangeStatus(ctx, id)
			if err != nil {
				return nil, err
			}
			pirp.State = "paid"
		}
	case "BTC":
		result, hash, err := cuc.InfoBTC(pirp)
		if err != nil {
			return nil, err
		}

		if result {
			err = cuc.repo.UpdateTransactionHash(ctx, hash, pirp.ID)
			if err != nil {
				return nil, err
			}

			err := cuc.repo.ChangeStatus(ctx, id)
			if err != nil {
				return nil, err
			}
			pirp.State = "paid"
		}
	}

	return pirp, nil
}

func (cuc *invoiceUseCase) InfoETH(pirp *models.PaymentInfoResponse) (bool, string, error) {
	var tm []*models.ETHTransaction

	url := fmt.Sprintf("https://api.ethplorer.io/getAddressTransactions/%s?apiKey=%s", cuc.cfg.ETHRecipient, cuc.cfg.Ethplorer)
	resp, err := http.Get(url)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, "", errors.New("api did not respond with a 200 code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return false, "", err
	}

	err = json.Unmarshal(body, &tm)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return false, "", err
	}

	for _, t := range tm {
		str := strconv.FormatFloat(t.Value, 'f', -1, 64)
		if t.From == pirp.FromAddress && str == pirp.Amount {
			result, err := cuc.repo.CheckTransactionHash(context.Background(), t.Hash)
			if err != nil {
				return false, "", nil
			}
			if !result {
				return true, t.Hash, nil
			}
			continue
		}
	}

	return false, "", nil
}

func (cuc *invoiceUseCase) InfoBTC(pirp *models.PaymentInfoResponse) (bool, string, error) {
	address, err := cuc.bcClient.GetAddress(cuc.cfg.BTCRecipient)
	if err != nil {
		return false, "", err
	}

	for _, t := range address.Txs {
		for _, i := range t.Inputs {
			if i.PrevOut.Addr == pirp.FromAddress {
				for _, o := range t.Out {
					btcValue := float64(o.Value) / 100000000.0
					strValue := strconv.FormatFloat(btcValue, 'f', -1, 64)
					if o.Addr == cuc.cfg.BTCRecipient && strValue == pirp.Amount {
						result, err := cuc.repo.CheckTransactionHash(context.Background(), t.Hash)
						if err != nil {
							return false, "", nil
						}
						if !result {
							return true, t.Hash, nil
						}
						continue
					}
				}
			}
		}
	}

	return false, "", nil
}

func (cuc *invoiceUseCase) CheckID(ctx context.Context, id string) (bool, error) {
	result, err := cuc.repo.CheckID(ctx, id)
	if err != nil {
		return false, err
	}

	return result, nil
}
