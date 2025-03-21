package model

import (
	"time"
)

// TransactionType は入庫または出庫を表す列挙型
type TransactionType string

const (
	TransactionTypeIn  TransactionType = "in"
	TransactionTypeOut TransactionType = "out"
)

// Transaction はトランザクションのドメインモデル
type Transaction struct {
	ID        string         `json:"id"`
	StockID   string         `json:"stock_id"`
	Type      TransactionType `json:"type"`
	Quantity  int            `json:"quantity"`
	Note      string         `json:"note"`
	CreatedAt time.Time      `json:"created_at"`
}

// NewTransaction は新しいトランザクションを作成する
func NewTransaction(stockID string, transType TransactionType, quantity int, note string) *Transaction {
	return &Transaction{
		StockID:   stockID,
		Type:      transType,
		Quantity:  quantity,
		Note:      note,
		CreatedAt: time.Now(),
	}
}

// GetAdjustmentQuantity は在庫調整量を返す（出庫の場合は負の値）
func (t *Transaction) GetAdjustmentQuantity() int {
	if t.Type == TransactionTypeOut {
		return -t.Quantity
	}
	return t.Quantity
}

// Validate はトランザクションが有効かどうかを検証する
func (t *Transaction) Validate() error {
	if t.StockID == "" {
		return ErrEmptyStockID
	}
	if t.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	if t.Type != TransactionTypeIn && t.Type != TransactionTypeOut {
		return ErrInvalidTransactionType
	}
	return nil
}
