package database

import (
	"service-3/domain/entities"
	"time"

	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) GetByID(id string) (*entities.Alert, error) {
	var model AlertModel

	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	return mapAlertModelToEntity(&model), nil
}

func (r *AlertRepository) List(isResolved *bool, stockID string) ([]entities.Alert, error) {
	var models []AlertModel
	db := r.db

	if isResolved != nil {
		db = db.Where("is_resolved = ?", *isResolved)
	}

	if stockID != "" {
		db = db.Where("stock_id = ?", stockID)
	}

	if result := db.Order("created_at desc").Find(&models); result.Error != nil {
		return nil, result.Error
	}

	// モデルからエンティティへマッピング
	alerts := make([]entities.Alert, len(models))
	for i, model := range models {
		alerts[i] = *mapAlertModelToEntity(&model)
	}

	return alerts, nil
}

func (r *AlertRepository) Create(alert *entities.Alert) error {
	model := mapAlertEntityToModel(alert)
	return r.db.Create(model).Error
}

func (r *AlertRepository) Update(alert *entities.Alert) error {
	model := mapAlertEntityToModel(alert)
	return r.db.Save(model).Error
}

func (r *AlertRepository) GetAlertSummary() (*entities.AlertSummary, error) {
	var summary entities.AlertSummary

	// 合計アラート数
	if err := r.db.Model(&AlertModel{}).Count(&summary.TotalAlerts).Error; err != nil {
		return nil, err
	}

	// 解決済みアラート数
	if err := r.db.Model(&AlertModel{}).Where("is_resolved = ?", true).Count(&summary.ResolvedAlerts).Error; err != nil {
		return nil, err
	}

	// アクティブなアラート数
	if err := r.db.Model(&AlertModel{}).Where("is_resolved = ?", false).Count(&summary.ActiveAlerts).Error; err != nil {
		return nil, err
	}

	// 最新のアラート日時
	var lastAlert AlertModel
	if err := r.db.Order("created_at desc").First(&lastAlert).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	summary.LastAlertDate = time.Now() // 仮の値としてのみ設定

	return &summary, nil
}

func mapAlertModelToEntity(model *AlertModel) *entities.Alert {
	return &entities.Alert{
		ID:         model.ID,
		StockID:    model.StockID,
		Type:       model.Type,
		Message:    model.Message,
		IsResolved: model.IsResolved,
		CreatedAt:  time.Now(), // 仮の値としてのみ設定
		ResolvedAt: nil,        // 実際のマッピングではここも適切な値に設定する
	}
}

func mapAlertEntityToModel(entity *entities.Alert) *AlertModel {
	return &AlertModel{
		ID:         entity.ID,
		StockID:    entity.StockID,
		Type:       entity.Type,
		Message:    entity.Message,
		IsResolved: entity.IsResolved,
		// 時間フィールドのマッピングは省略
	}
}
