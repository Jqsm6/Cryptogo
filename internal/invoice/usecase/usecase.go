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

func (cuc *invoiceUseCase) Create(ctx context.Context, paymentRequest *models.PaymentRequest) (*models.PaymentResponse, error) {
	var paymentResponse models.PaymentResponse

	switch paymentRequest.Currency {
	case "ETH":
		paymentResponse.ToAddress = cuc.cfg.ETH
	case "BTC":
		paymentResponse.ToAddress = cuc.cfg.BTC
	case "BNB":
		paymentResponse.ToAddress = cuc.cfg.BNB
	}

	guid := xid.New().String()
	paymentResponse.ID = guid

	amount := strconv.FormatFloat(paymentRequest.Amount, 'f', -1, 64)
	paymentDB := models.PaymentDB{
		ID:          guid,
		State:       "notpayed",
		Currency:    paymentRequest.Currency,
		Amount:      amount,
		ToAddress:   paymentResponse.ToAddress,
		FromAddress: paymentRequest.FromAddress,
	}

	err := cuc.repo.Create(ctx, &paymentDB)
	if err != nil {
		return nil, err
	}

	return &paymentResponse, nil
}

func (cuc *invoiceUseCase) Info(ctx context.Context, paymentInfoRequest *models.PaymentInfoRequest) (*models.PaymentInfoResponse, error) {
	id := paymentInfoRequest.ID

	paymentInfoResponse, err := cuc.repo.Info(ctx, id)
	if err != nil {
		return nil, err
	}

	if paymentInfoResponse.State == "paid" {
		return paymentInfoResponse, nil
	}

	switch paymentInfoResponse.Currency {
	case "ETH":
		result, hash, err := cuc.InfoETH(paymentInfoResponse)
		if err != nil {
			return nil, err
		}

		if result {
			err = cuc.repo.UpdateHash(ctx, hash, paymentInfoResponse.ID)
			if err != nil {
				return nil, err
			}

			err := cuc.repo.ChangeStatus(ctx, id)
			if err != nil {
				return nil, err
			}
			paymentInfoResponse.State = "paid"
		}
	case "BTC":
		result, hash, err := cuc.InfoBTC(paymentInfoResponse)
		if err != nil {
			return nil, err
		}

		if result {
			err = cuc.repo.UpdateHash(ctx, hash, paymentInfoResponse.ID)
			if err != nil {
				return nil, err
			}

			err := cuc.repo.ChangeStatus(ctx, id)
			if err != nil {
				return nil, err
			}
			paymentInfoResponse.State = "paid"
		}
	case "BNB":
		result, hash, err := cuc.InfoBNB(paymentInfoResponse)
		if err != nil {
			return nil, err
		}

		if result {
			err = cuc.repo.UpdateHash(ctx, hash, paymentInfoResponse.ID)
			if err != nil {
				return nil, err
			}

			err := cuc.repo.ChangeStatus(ctx, id)
			if err != nil {
				return nil, err
			}
			paymentInfoResponse.State = "paid"
		}
	}

	return paymentInfoResponse, nil
}

func (cuc *invoiceUseCase) InfoETH(paymentInfoResponse *models.PaymentInfoResponse) (bool, string, error) {
	var ethTransactions []*models.ETHTransaction

	url := fmt.Sprintf("https://api.ethplorer.io/getAddressTransactions/%s?apiKey=%s", cuc.cfg.ETH, cuc.cfg.Ethplorer)
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

	err = json.Unmarshal(body, &ethTransactions)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return false, "", err
	}

	for _, tx := range ethTransactions {
		amount := strconv.FormatFloat(tx.Value, 'f', -1, 64)
		if tx.From == paymentInfoResponse.FromAddress && amount == paymentInfoResponse.Amount {
			result, err := cuc.repo.CheckHash(context.Background(), tx.Hash)
			if err != nil {
				return false, "", nil
			}
			if !result {
				return true, tx.Hash, nil
			}
			continue
		}
	}

	return false, "", nil
}

func (cuc *invoiceUseCase) InfoBTC(paymentInfoResponse *models.PaymentInfoResponse) (bool, string, error) {
	address, err := cuc.bcClient.GetAddress(cuc.cfg.BTC)
	if err != nil {
		return false, "", err
	}

	for _, tx := range address.Txs {
		for _, i := range tx.Inputs {
			if i.PrevOut.Addr == paymentInfoResponse.FromAddress {
				for _, o := range tx.Out {
					amount := float64(o.Value) / 100000000.0
					stringAmount := strconv.FormatFloat(amount, 'f', -1, 64)
					if o.Addr == cuc.cfg.BTC && stringAmount == paymentInfoResponse.Amount {
						result, err := cuc.repo.CheckHash(context.Background(), tx.Hash)
						if err != nil {
							return false, "", nil
						}
						if !result {
							return true, tx.Hash, nil
						}
						continue
					}
				}
			}
		}
	}

	return false, "", nil
}

func (cuc *invoiceUseCase) InfoBNB(paymentInfoResponse *models.PaymentInfoResponse) (bool, string, error) {
	var address *models.BNBTransaction

	url := fmt.Sprintf("https://api.bscscan.com/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=50&sort=asc&apikey=%s", cuc.cfg.BNB, cuc.cfg.Bscscan)
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

	err = json.Unmarshal(body, &address)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return false, "", err
	}

	for _, tx := range address.Result {
		amount, err := strconv.Atoi(tx.Value)
		if err != nil {
			return false, "", err
		}
		floatAmount := float64(amount) / 1000000000000000000.0
		stringAmount := strconv.FormatFloat(floatAmount, 'f', -1, 64)
		if tx.From == paymentInfoResponse.FromAddress && stringAmount == paymentInfoResponse.Amount {
			result, err := cuc.repo.CheckHash(context.Background(), tx.Hash)
			if err != nil {
				return false, "", nil
			}
			if !result {
				return true, tx.Hash, nil
			}
			continue
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
