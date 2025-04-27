package repositories

import (
	"service-3/domain/entities"
)

type AlertRepository interface {
	GetByID(id string) (*entities.Alert, error)
	List(isResolved *bool, stockID string) ([]entities.Alert, error)
	Create(alert *entities.Alert) error
	Update(alert *entities.Alert) error
	GetAlertSummary() (*entities.AlertSummary, error)
}
