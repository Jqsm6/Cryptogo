package usecase

import (
	"net/http"

	"Cryptogo/internal/models"
	"Cryptogo/internal/status"
	"Cryptogo/pkg/logger"
)

type statusUseCase struct {
	log *logger.Logger
}

func NewStatusUseCase(log *logger.Logger) status.UseCase {
	return &statusUseCase{log: log}
}

func (suc *statusUseCase) GetAPIStatus() *models.Status {
	asModel := models.Status{}

	ethResp, err := http.Get(ethAPI)
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
