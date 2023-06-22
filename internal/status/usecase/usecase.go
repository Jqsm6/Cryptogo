package usecase

import (
	"fmt"
	"net/http"

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

func (suc *statusUseCase) GetAPIStatus() *models.Status {
	asModel := models.Status{}

	ethUrl := fmt.Sprintf("https://api.ethplorer.io/getTokenInfo/0xf3db5fa2c66b7af3eb0c0b782510816cbe4813b8?apiKey=%s", suc.cfg.Ethplorer)
	ethResp, err := http.Get(ethUrl)
	if err != nil {
		suc.log.Info().Err(err).Msg("usecase")
	}
	defer ethResp.Body.Close()

	if ethResp.StatusCode == 200 {
		asModel.API.Status.ETH = "ok"
	} else {
		asModel.API.Status.ETH = "bad"
	}

	return &asModel
}
