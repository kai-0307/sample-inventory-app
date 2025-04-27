package database

import (
	"fmt"

	"service-3/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// データベースのモデル定義（GORMのためのもの）
type AlertModel struct {
	ID         string `gorm:"primaryKey"`
	StockID    string `json:"stock_id"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	IsResolved bool   `json:"is_resolved"`
	CreatedAt  gorm.DeletedAt
	ResolvedAt *gorm.DeletedAt
}

type AlertConfigModel struct {
	ID          string `gorm:"primaryKey"`
	StockID     string `json:"stock_id"`
	MinQuantity int    `json:"min_quantity"`
	MaxQuantity int    `json:"max_quantity"`
	IsActive    bool   `json:"is_active"`
	UpdatedAt   gorm.DeletedAt
}

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// マイグレーション実行
	db.AutoMigrate(&AlertModel{}, &AlertConfigModel{})

	return db, nil
}
