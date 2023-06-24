package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/xid"

	"Cryptogo/config"
	"Cryptogo/internal/invoice"
	"Cryptogo/internal/models"
	"Cryptogo/pkg/logger"
)

type invoiceUseCase struct {
	repo invoice.Repository
	log  *logger.Logger
	cfg  *config.Config
}

func NewInvoiceUseCase(repo invoice.Repository, log *logger.Logger, cfg *config.Config) invoice.UseCase {
	return &invoiceUseCase{repo: repo, log: log, cfg: cfg}
}

func (cuc *invoiceUseCase) Create(ctx context.Context, paymentRequest *models.PaymentRequest) (*models.PaymentResponse, error) {
	var paymentResponse models.PaymentResponse

	switch paymentRequest.Currency {
	case "ETH":
		paymentResponse.ToAddress = paymentRequest.ToAddress
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

	return paymentInfoResponse, nil
}

func (cuc *invoiceUseCase) ConfirmETH(ctx context.Context, paymentConfirmRequest *models.PaymentConfirmRequest) (*models.PaymentConfirmResponse, error) {
	var ethTransaction *models.ETHTransaction
	var paymentConfirmResponse models.PaymentConfirmResponse
	result, err := cuc.repo.CheckHash(ctx, paymentConfirmRequest.TxHash)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return nil, err
	}

	if result {
		return nil, errors.New("was earlier")
	}

	node := fmt.Sprintf("https://eth.getblock.io/%s/mainnet/", cuc.cfg.Ethereum)
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionByHash",
		"params": []interface{}{
			paymentConfirmRequest.TxHash,
		},
		"id": 1,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return nil, err
	}

	req, err := http.NewRequest("POST", node, bytes.NewBuffer(jsonData))
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return nil, err
	}

	err = json.Unmarshal(responseData, &ethTransaction)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return nil, err
	}

	infoDB, err := cuc.repo.Get(ctx, paymentConfirmRequest.ID)
	if err != nil {
		return nil, err
	}

	stringAmount := ethTransaction.Result.Value[2:]
	intAmount, err := strconv.ParseInt(stringAmount, 16, 64)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return nil, err
	}
	floatAmount := float64(intAmount) / 1000000000000000000.0
	stringAmount = strconv.FormatFloat(floatAmount, 'f', -1, 64)

	if ethTransaction.Result.From == strings.ToLower(infoDB.FromAddress) && ethTransaction.Result.To == strings.ToLower(infoDB.ToAddress) && stringAmount == infoDB.Amount {
		paymentConfirmResponse.State = "found"

		err := cuc.repo.ChangeStatus(ctx, paymentConfirmRequest.ID)
		if err != nil {
			cuc.log.Err(err).Msg("usecase")
			return nil, err
		}

		err = cuc.repo.UpdateHash(ctx, paymentConfirmRequest.TxHash, paymentConfirmRequest.ID)
		if err != nil {
			cuc.log.Err(err).Msg("usecase")
			return nil, err
		}

		return &paymentConfirmResponse, nil
	}
	paymentConfirmResponse.State = "notfound"

	return &paymentConfirmResponse, nil
}

func (cuc *invoiceUseCase) CheckID(ctx context.Context, id string) (bool, error) {
	result, err := cuc.repo.CheckID(ctx, id)
	if err != nil {
		return false, err
	}

	return result, nil
}
