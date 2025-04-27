package usecases

import (
	"service-3/domain/entities"
	"service-3/domain/repositories"
	"service-3/domain/services"
)

type ReportUseCase struct {
	alertRepo    repositories.AlertRepository
	configRepo   repositories.AlertConfigRepository
	stockService services.StockService
}

func NewReportUseCase(
	alertRepo repositories.AlertRepository,
	configRepo repositories.AlertConfigRepository,
	stockService services.StockService,
) *ReportUseCase {
	return &ReportUseCase{
		alertRepo:    alertRepo,
		configRepo:   configRepo,
		stockService: stockService,
	}
}

func (uc *ReportUseCase) GenerateStockReport() ([]entities.StockReport, error) {
	reports, err := uc.configRepo.GetStockReportData()
	if err != nil {
		return nil, err
	}

	// 現在の在庫数を取得
	for i := range reports {
		currentStock, err := uc.stockService.GetCurrentStock(reports[i].StockID)
		if err != nil {
			// エラーハンドリング（この例ではエラーを無視して続行）
			continue
		}
		reports[i].CurrentStock = currentStock
	}

	return reports, nil
}

func (uc *ReportUseCase) GenerateAlertReport() (*entities.AlertSummary, error) {
	return uc.alertRepo.GetAlertSummary()
}
