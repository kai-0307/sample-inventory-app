package database

import (
	"service-3/domain/entities"
	"time"

	"gorm.io/gorm"
)

type AlertConfigRepository struct {
	db *gorm.DB
}

func NewAlertConfigRepository(db *gorm.DB) *AlertConfigRepository {
	return &AlertConfigRepository{db: db}
}

func (r *AlertConfigRepository) GetByStockID(stockID string) (*entities.AlertConfig, error) {
	var model AlertConfigModel

	if err := r.db.Where("stock_id = ?", stockID).First(&model).Error; err != nil {
		return nil, err
	}

	return mapAlertConfigModelToEntity(&model), nil
}

func (r *AlertConfigRepository) Create(config *entities.AlertConfig) error {
	model := mapAlertConfigEntityToModel(config)
	return r.db.Create(model).Error
}

func (r *AlertConfigRepository) Update(config *entities.AlertConfig) error {
	model := mapAlertConfigEntityToModel(config)
	return r.db.Save(model).Error
}

func (r *AlertConfigRepository) GetStockReportData() ([]entities.StockReport, error) {
	var reports []entities.StockReport

	rows, err := r.db.Raw(`
      SELECT
          acm.stock_id,
          acm.min_quantity,
          acm.max_quantity,
          COUNT(am.id) as alert_count,
          MAX(am.created_at) as last_alert
      FROM
          alert_config_models acm
          LEFT JOIN alert_models am ON acm.stock_id = am.stock_id
      WHERE
          acm.is_active = true
      GROUP BY
          acm.stock_id, acm.min_quantity, acm.max_quantity
  `).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report entities.StockReport
		if err := rows.Scan(
			&report.StockID,
			&report.MinQuantity,
			&report.MaxQuantity,
			&report.AlertCount,
			&report.LastAlert,
		); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func mapAlertConfigModelToEntity(model *AlertConfigModel) *entities.AlertConfig {
	return &entities.AlertConfig{
		ID:          model.ID,
		StockID:     model.StockID,
		MinQuantity: model.MinQuantity,
		MaxQuantity: model.MaxQuantity,
		IsActive:    model.IsActive,
		UpdatedAt:   time.Now(), // 仮の値としてのみ設定
	}
}

func mapAlertConfigEntityToModel(entity *entities.AlertConfig) *AlertConfigModel {
	return &AlertConfigModel{
		ID:          entity.ID,
		StockID:     entity.StockID,
		MinQuantity: entity.MinQuantity,
		MaxQuantity: entity.MaxQuantity,
		IsActive:    entity.IsActive,
		// 時間フィールドのマッピングは省略
	}
}
