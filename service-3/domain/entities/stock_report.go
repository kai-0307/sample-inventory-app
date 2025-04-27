package entities

import (
	"time"
)

type StockReport struct {
	StockID      string    `json:"stock_id"`
	CurrentStock int       `json:"current_stock"`
	MinQuantity  int       `json:"min_quantity"`
	MaxQuantity  int       `json:"max_quantity"`
	AlertCount   int       `json:"alert_count"`
	LastAlert    time.Time `json:"last_alert"`
}
