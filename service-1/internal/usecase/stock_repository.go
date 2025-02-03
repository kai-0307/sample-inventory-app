package usecase

import (
	"context"

	"service-1/internal/domain/models"
)

type StockRepository interface {
    Create(ctx context.Context, stock *models.Stock) error
    FindByID(ctx context.Context, id string) (*models.Stock, error)
    FindAll(ctx context.Context) ([]*models.Stock, error)
    FindLowStock(ctx context.Context) ([]*models.Stock, error)
    Update(ctx context.Context, stock *models.Stock) error
    Delete(ctx context.Context, id string) error
    UpdateQuantity(ctx context.Context, id string, adjustment int) error
}
