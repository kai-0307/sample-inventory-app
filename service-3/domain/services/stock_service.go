package services

type StockService interface {
	GetCurrentStock(stockID string) (int, error)
}
