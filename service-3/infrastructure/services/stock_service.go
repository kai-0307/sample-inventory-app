package services

import (
	"service-3/domain/services"
)

type StockServiceImpl struct {
	// 必要に応じて他のサービス（HTTP クライアントなど）を注入
}

func NewStockService() services.StockService {
	return &StockServiceImpl{}
}

func (s *StockServiceImpl) GetCurrentStock(stockID string) (int, error) {
	// 実際の実装ではHTTPリクエストでStock Serviceから取得
	// ここではダミー値を返す
	return 100, nil
}
