package entities

import (
	"time"
)

type AlertConfig struct {
	ID          string    `json:"id"`
	StockID     string    `json:"stock_id"`
	MinQuantity int       `json:"min_quantity"`
	MaxQuantity int       `json:"max_quantity"`
	IsActive    bool      `json:"is_active"`
	UpdatedAt   time.Time `json:"updated_at"`
}
