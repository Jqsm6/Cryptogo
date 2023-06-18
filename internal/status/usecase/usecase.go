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

func (suc *statusUseCase) GetAPIStatus() *models.APIStatus {
	asModel := models.APIStatus{}

	ethResp, err := http.Get(ethURL)
	if err != nil {
		suc.log.Err(err)
	}
	defer ethResp.Body.Close()

	if ethResp.StatusCode == 200 {
		asModel.StatusETH = "ok"
	} else {
		asModel.StatusETH = "bad"
	}

	return &asModel
}
