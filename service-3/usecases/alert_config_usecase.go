package usecases

import (
	"time"

	"service-3/domain/entities"
	"service-3/domain/repositories"
)

type AlertConfigUseCase struct {
	repo repositories.AlertConfigRepository
}

func NewAlertConfigUseCase(repo repositories.AlertConfigRepository) *AlertConfigUseCase {
	return &AlertConfigUseCase{repo: repo}
}

func (uc *AlertConfigUseCase) GetAlertConfig(stockID string) (*entities.AlertConfig, error) {
	return uc.repo.GetByStockID(stockID)
}

func (uc *AlertConfigUseCase) CreateAlertConfig(config *entities.AlertConfig) error {
	config.UpdatedAt = time.Now()
	config.IsActive = true
	return uc.repo.Create(config)
}

func (uc *AlertConfigUseCase) UpdateAlertConfig(config *entities.AlertConfig) error {
	config.UpdatedAt = time.Now()
	return uc.repo.Update(config)
}
