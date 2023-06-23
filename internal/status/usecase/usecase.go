package usecase

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"

	"Cryptogo/config"
	"Cryptogo/internal/models"
	"Cryptogo/internal/status"
	"Cryptogo/pkg/logger"
)

type statusUseCase struct {
	log *logger.Logger
	cfg *config.Config
}

func NewStatusUseCase(log *logger.Logger, cfg *config.Config) status.UseCase {
	return &statusUseCase{log: log, cfg: cfg}
}

func (suc *statusUseCase) GetAPIStatus() (*models.Status, error) {
	asModel := models.Status{}

	ethUrl := fmt.Sprintf("https://api.ethplorer.io/getTokenInfo/0xf3db5fa2c66b7af3eb0c0b782510816cbe4813b8?apiKey=%s", suc.cfg.Ethplorer)
	ethResp, err := http.Get(ethUrl)
	if err != nil {
		suc.log.Err(err).Msg("usecase")
		return nil, err
	}
	defer ethResp.Body.Close()

	if ethResp.StatusCode == 200 {
		asModel.API.Status.ETH = "ok"
	} else {
		asModel.API.Status.ETH = "bad"
	}

	btcUrl := "https://blockchain.info/latestblock"
	btcResp, err := http.Get(btcUrl)
	if err != nil {
		suc.log.Err(err).Msg("usecase")
		return nil, err
	}
	defer btcResp.Body.Close()

	if btcResp.StatusCode == 200 {
		asModel.API.Status.BTC = "ok"
	} else {
		asModel.API.Status.BTC = "bad"
	}

	var bsr *models.BNBStatusResponse
	bnbUrl := fmt.Sprintf("https://api.bscscan.com/api?module=gastracker&action=gasoracle&apikey=%s", suc.cfg.Bscscan)
	bnbResp, err := http.Get(bnbUrl)
	if err != nil {
		suc.log.Err(err).Msg("usecase")
		return nil, err
	}
	defer bnbResp.Body.Close()

	body, err := io.ReadAll(bnbResp.Body)
	if err != nil {
		suc.log.Err(err).Msg("usecase")
		return nil, err
	}

	err = json.Unmarshal(body, &bsr)
	if err != nil {
		suc.log.Err(err).Msg("usecase")
		return nil, err
	}

	if bsr.Status == "1" {
		asModel.API.Status.BNB = "ok"
	} else {
		asModel.API.Status.BNB = "bad"
	}

	return &asModel, nil
}
