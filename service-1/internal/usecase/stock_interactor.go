package usecase

import (
	"context"

	"service-1/internal/domain/models"
)

type StockInteractor struct {
    repository StockRepository
}

func NewStockInteractor(repository StockRepository) *StockInteractor {
    return &StockInteractor{
        repository: repository,
    }
}

func (i *StockInteractor) CreateStock(ctx context.Context, input *models.Stock) error {
    stock, err := models.NewStock(
        input.ID,
        input.Name,
        input.Description,
        input.Quantity,
        input.Unit,
        input.MinQuantity,
    )
    if err != nil {
        return err
    }

    return i.repository.Create(ctx, stock)
}

func (i *StockInteractor) GetStock(ctx context.Context, id string) (*models.Stock, error) {
    return i.repository.FindByID(ctx, id)
}

func (i *StockInteractor) ListStocks(ctx context.Context) ([]*models.Stock, error) {
    return i.repository.FindAll(ctx)
}

func (i *StockInteractor) ListLowStocks(ctx context.Context) ([]*models.Stock, error) {
    return i.repository.FindLowStock(ctx)
}

func (i *StockInteractor) UpdateStock(ctx context.Context, stock *models.Stock) error {
    return i.repository.Update(ctx, stock)
}

func (i *StockInteractor) DeleteStock(ctx context.Context, id string) error {
    return i.repository.Delete(ctx, id)
}

func (i *StockInteractor) UpdateStockQuantity(ctx context.Context, id string, adjustment int) error {
    stock, err := i.repository.FindByID(ctx, id)
    if err != nil {
        return err
    }

    if err := stock.UpdateQuantity(adjustment); err != nil {
        return err
    }

    return i.repository.Update(ctx, stock)
}
