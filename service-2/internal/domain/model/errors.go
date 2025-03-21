package model

import "errors"

var (
	ErrEmptyStockID          = errors.New("stock ID cannot be empty")
	ErrInvalidQuantity       = errors.New("quantity must be greater than 0")
	ErrInvalidTransactionType = errors.New("transaction type must be either 'in' or 'out'")
	ErrTransactionNotFound   = errors.New("transaction not found")
	ErrStockNotFound         = errors.New("stock not found")
	ErrDatabaseError         = errors.New("database error")
	ErrExternalServiceError  = errors.New("external service error")
)
