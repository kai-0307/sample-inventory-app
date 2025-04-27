package repositories

import (
	"service-3/domain/entities"
)

type AlertConfigRepository interface {
	GetByStockID(stockID string) (*entities.AlertConfig, error)
	Create(config *entities.AlertConfig) error
	Update(config *entities.AlertConfig) error
	GetStockReportData() ([]entities.StockReport, error)
}
