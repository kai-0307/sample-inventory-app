package entities

import (
	"time"
)

type Alert struct {
	ID         string     `json:"id"`
	StockID    string     `json:"stock_id"`
	Type       string     `json:"type"` // "low_stock", "excess_stock", etc.
	Message    string     `json:"message"`
	IsResolved bool       `json:"is_resolved"`
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at"`
}
