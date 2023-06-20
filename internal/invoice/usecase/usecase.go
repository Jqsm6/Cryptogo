package usecase

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (cuc *invoiceUseCase) Create(ctx context.Context, prqm *models.PaymentRequest) (*models.PaymentResponse, error) {
	var prpm models.PaymentResponse

	guid := xid.New().String()

	var privateKeyBytes []byte
	var address string

	switch prqm.Currency {
	case "ETH":
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			cuc.log.Err(err).Msg("usecase")
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
		ID:         guid,
		State:      "notpayed",
		Currency:   prqm.Currency,
		Amount:     prqm.Amount,
		ToAddress:  address,
		PrivateKey: hexutil.Encode(privateKeyBytes),
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
		result, err := cuc.InfoETH(pirp)
		if err != nil {
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

func (cuc *invoiceUseCase) InfoETH(pirp *models.PaymentInfoResponse) (bool, error) {
	var em *models.ETHBalance

	url := fmt.Sprintf("https://api.ethplorer.io/getAddressInfo/%s?apiKey=freekey", pirp.ToAddress)
	resp, err := http.Get(url)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, errors.New("api did not respond with a 200 code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return false, err
	}

	err = json.Unmarshal(body, &em)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return false, err
	}

	balance := strconv.Itoa(em.ETH.Balance)

	if balance == pirp.Amount {
		err = cuc.WithdrawETH(pirp)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (cuc *invoiceUseCase) CheckID(ctx context.Context, id string) (bool, error) {
	result, err := cuc.repo.CheckID(ctx, id)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (cuc *invoiceUseCase) WithdrawETH(pirp *models.PaymentInfoResponse) error {
	key, err := cuc.repo.GetPrivateKey(pirp.ID)
	if err != nil {
		return err
	}
	key = key[2:]

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		println(1)
		cuc.log.Err(err).Msg("usecase")
		return err
	}

	client, err := ethclient.Dial(fmt.Sprintf("https://eth-mainnet.g.alchemy.com/v2/%s", cuc.cfg.AlchemyKey))
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		cuc.log.Err(err).Msg("usecase")
		return err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return err
	}

	toAddress := common.HexToAddress(pirp.ToAddress)
	value, _ := new(big.Int).SetString(pirp.Amount, 16)

	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return err
	}

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		cuc.log.Err(err).Msg("usecase")
		return err
	}

	cuc.log.Info().Msgf("Transaction is complete. Hash: %s", signedTx.Hash().Hex())
	return nil
}
