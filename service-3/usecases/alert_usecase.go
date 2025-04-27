package usecases

import (
	"time"

	"service-3/domain/entities"
	"service-3/domain/repositories"
)

type AlertUseCase struct {
	repo repositories.AlertRepository
}

func NewAlertUseCase(repo repositories.AlertRepository) *AlertUseCase {
	return &AlertUseCase{repo: repo}
}

func (uc *AlertUseCase) GetAlert(id string) (*entities.Alert, error) {
	return uc.repo.GetByID(id)
}

func (uc *AlertUseCase) ListAlerts(isResolved *bool, stockID string) ([]entities.Alert, error) {
	return uc.repo.List(isResolved, stockID)
}

func (uc *AlertUseCase) ResolveAlert(id string) (*entities.Alert, error) {
	alert, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	alert.IsResolved = true
	alert.ResolvedAt = &now

	if err := uc.repo.Update(alert); err != nil {
		return nil, err
	}

	return alert, nil
}
