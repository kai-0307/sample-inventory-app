package service

import (
	"service-2/internal/domain/model"
	"service-2/internal/repository"
)

// ITransactionService はトランザクションサービスのインターフェース
type ITransactionService interface {
	CreateTransaction(stockID string, transType string, quantity int, note string) (*model.Transaction, error)
	GetTransaction(id string) (*model.Transaction, error)
	ListTransactions(limit int, startDate, endDate string) ([]*model.Transaction, error)
	GetStockTransactions(stockID string) ([]*model.Transaction, error)
	GetStockSummary(stockID string) (*model.TransactionSummary, error)
}

// TransactionService はトランザクションサービスの実装
type TransactionService struct {
	transactionRepo repository.ITransactionRepository
	stockRepo       repository.IStockRepository
}

// NewTransactionService は新しいトランザクションサービスを作成する
func NewTransactionService(
	transactionRepo repository.ITransactionRepository,
	stockRepo repository.IStockRepository,
) ITransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		stockRepo:       stockRepo,
	}
}

// CreateTransaction は新しいトランザクションを作成する
func (s *TransactionService) CreateTransaction(stockID string, transType string, quantity int, note string) (*model.Transaction, error) {
	// トランザクションモデルの作成
	transaction := model.NewTransaction(stockID, model.TransactionType(transType), quantity, note)

	// バリデーション
	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	// トランザクションの保存
	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	// 在庫数の更新
	adjustmentQty := transaction.GetAdjustmentQuantity()
	if err := s.stockRepo.UpdateStockQuantity(stockID, adjustmentQty, "Transaction update"); err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransaction はIDによってトランザクションを取得する
func (s *TransactionService) GetTransaction(id string) (*model.Transaction, error) {
	return s.transactionRepo.FindByID(id)
}

// ListTransactions はトランザクション一覧を取得する
func (s *TransactionService) ListTransactions(limit int, startDate, endDate string) ([]*model.Transaction, error) {
	if limit <= 0 {
		limit = 100 // デフォルト値
	}
	return s.transactionRepo.FindAll(limit, startDate, endDate)
}

// GetStockTransactions は在庫IDによってトランザクションを取得する
func (s *TransactionService) GetStockTransactions(stockID string) ([]*model.Transaction, error) {
	return s.transactionRepo.FindByStockID(stockID)
}

// GetStockSummary は在庫の集計情報を取得する
func (s *TransactionService) GetStockSummary(stockID string) (*model.TransactionSummary, error) {
	// トランザクション集計の取得
	summary, err := s.transactionRepo.GetStockTransactionSummary(stockID)
	if err != nil {
		return nil, err
	}

	// 在庫情報の取得
	stock, err := s.stockRepo.GetStock(stockID)
	if err != nil {
		return nil, err
	}

	// 在庫残高の設定
	summary.Balance = stock.Quantity

	return summary, nil
}


