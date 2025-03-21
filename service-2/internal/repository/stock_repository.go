package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"service-2/internal/domain/model"
)

// IStockRepository は在庫リポジトリのインターフェース
type IStockRepository interface {
	GetStock(stockID string) (*model.Stock, error)
	UpdateStockQuantity(stockID string, adjustment int, note string) error
}

// StockRepository は在庫リポジトリの実装
type StockRepository struct {
	baseURL string
	client  *http.Client
}

// NewStockRepository は新しい在庫リポジトリを作成する
func NewStockRepository(baseURL string) IStockRepository {
	return &StockRepository{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// GetStock は在庫情報を取得する
func (r *StockRepository) GetStock(stockID string) (*model.Stock, error) {
	url := fmt.Sprintf("%s/api/service-1/stocks/%s", r.baseURL, stockID)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, model.ErrStockNotFound
	}

	var stock model.Stock
	if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
		return nil, err
	}

	return &stock, nil
}

// UpdateStockQuantity は在庫数を更新する
func (r *StockRepository) UpdateStockQuantity(stockID string, adjustment int, note string) error {
	url := fmt.Sprintf("%s/api/service-1/stocks/%s/quantity", r.baseURL, stockID)

	payload := map[string]interface{}{
		"adjustment": adjustment,
		"note":       note,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.ErrExternalServiceError
	}

	return nil
}
