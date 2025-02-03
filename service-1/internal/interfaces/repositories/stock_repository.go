package repositories

import (
	"context"

	"service-1/internal/domain/models"

	"gorm.io/gorm"
)

type PostgresStockRepository struct {
    db *gorm.DB
}

func NewPostgresStockRepository(db *gorm.DB) *PostgresStockRepository {
    return &PostgresStockRepository{
        db: db,
    }
}

func (r *PostgresStockRepository) Create(ctx context.Context, stock *models.Stock) error {
    return r.db.WithContext(ctx).Create(stock).Error
}

func (r *PostgresStockRepository) FindByID(ctx context.Context, id string) (*models.Stock, error) {
    var stock models.Stock
    if err := r.db.WithContext(ctx).First(&stock, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &stock, nil
}

func (r *PostgresStockRepository) FindAll(ctx context.Context) ([]*models.Stock, error) {
    var stocks []*models.Stock
    if err := r.db.WithContext(ctx).Find(&stocks).Error; err != nil {
        return nil, err
    }
    return stocks, nil
}

func (r *PostgresStockRepository) FindLowStock(ctx context.Context) ([]*models.Stock, error) {
    var stocks []*models.Stock
    if err := r.db.WithContext(ctx).Where("quantity <= min_quantity").Find(&stocks).Error; err != nil {
        return nil, err
    }
    return stocks, nil
}

func (r *PostgresStockRepository) Update(ctx context.Context, stock *models.Stock) error {
    return r.db.WithContext(ctx).Save(stock).Error
}

func (r *PostgresStockRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Delete(&models.Stock{}, "id = ?", id).Error
}

func (r *PostgresStockRepository) UpdateQuantity(ctx context.Context, id string, adjustment int) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var stock models.Stock
        if err := tx.WithContext(ctx).First(&stock, "id = ?", id).Error; err != nil {
            return err
        }

        if err := stock.UpdateQuantity(adjustment); err != nil {
            return err
        }

        return tx.WithContext(ctx).Save(&stock).Error
    })
}
