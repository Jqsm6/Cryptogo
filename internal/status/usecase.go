package status

import "Cryptogo/internal/models"

type UseCase interface {
	GetAPIStatus() *models.APIStatus
}
