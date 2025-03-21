package repository

import (
	"time"

	"gorm.io/gorm"

	"service-2/internal/domain/model"
)

// TransactionEntity はデータベース内のトランザクションエンティティ
type TransactionEntity struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StockID   string    `gorm:"index"`
	Type      string
	Quantity  int
	Note      string
	CreatedAt time.Time
}

// ToModel はエンティティをドメインモデルに変換する
func (e *TransactionEntity) ToModel() *model.Transaction {
	return &model.Transaction{
		ID:        e.ID,
		StockID:   e.StockID,
		Type:      model.TransactionType(e.Type),
		Quantity:  e.Quantity,
		Note:      e.Note,
		CreatedAt: e.CreatedAt,
	}
}

// FromModel はドメインモデルをエンティティに変換する
func (e *TransactionEntity) FromModel(m *model.Transaction) {
	e.ID = m.ID
	e.StockID = m.StockID
	e.Type = string(m.Type)
	e.Quantity = m.Quantity
	e.Note = m.Note
	e.CreatedAt = m.CreatedAt
}

// TransactionRepository はトランザクションリポジトリのインターフェース
type ITransactionRepository interface {
	Create(transaction *model.Transaction) error
	FindByID(id string) (*model.Transaction, error)
	FindAll(limit int, startDate, endDate string) ([]*model.Transaction, error)
	FindByStockID(stockID string) ([]*model.Transaction, error)
	GetStockTransactionSummary(stockID string) (*model.TransactionSummary, error)
}

// TransactionRepository はトランザクションリポジトリの実装
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository は新しいトランザクションリポジトリを作成する
func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionRepository{db: db}
}

// Create はトランザクションを作成する
func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	entity := &TransactionEntity{}
	entity.FromModel(transaction)

	if err := r.db.Create(entity).Error; err != nil {
		return err
	}

	transaction.ID = entity.ID
	return nil
}

// FindByID はIDによってトランザクションを検索する
func (r *TransactionRepository) FindByID(id string) (*model.Transaction, error) {
	var entity TransactionEntity
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrTransactionNotFound
		}
		return nil, err
	}
	return entity.ToModel(), nil
}

// FindAll はすべてのトランザクションを取得する
func (r *TransactionRepository) FindAll(limit int, startDate, endDate string) ([]*model.Transaction, error) {
	var entities []TransactionEntity

	query := r.db
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Order("created_at desc").Limit(limit).Find(&entities).Error; err != nil {
		return nil, err
	}

	transactions := make([]*model.Transaction, len(entities))
	for i, entity := range entities {
		transactions[i] = entity.ToModel()
	}

	return transactions, nil
}

// FindByStockID は在庫IDによってトランザクションを検索する
func (r *TransactionRepository) FindByStockID(stockID string) ([]*model.Transaction, error) {
	var entities []TransactionEntity

	if err := r.db.Where("stock_id = ?", stockID).Order("created_at desc").Find(&entities).Error; err != nil {
		return nil, err
	}

	transactions := make([]*model.Transaction, len(entities))
	for i, entity := range entities {
		transactions[i] = entity.ToModel()
	}

	return transactions, nil
}

// GetStockTransactionSummary は在庫のトランザクション集計を取得する
func (r *TransactionRepository) GetStockTransactionSummary(stockID string) (*model.TransactionSummary, error) {
	summary := &model.TransactionSummary{
		StockID: stockID,
	}

	// 入庫合計の計算
	if err := r.db.Model(&TransactionEntity{}).
		Where("stock_id = ? AND type = ?", stockID, "in").
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&summary.TotalIn).Error; err != nil {
		return nil, err
	}

	// 出庫合計の計算
	if err := r.db.Model(&TransactionEntity{}).
		Where("stock_id = ? AND type = ?", stockID, "out").
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&summary.TotalOut).Error; err != nil {
		return nil, err
	}

	return summary, nil
}

// Migrate はデータベースマイグレーションを実行する
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&TransactionEntity{})
}
