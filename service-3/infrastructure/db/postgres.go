package database

import (
	"fmt"
	"log"
	"time"

	"service-3/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// データベースのモデル定義（GORMのためのもの）
type AlertModel struct {
	ID         string `gorm:"primaryKey"`
	StockID    string
	Type       string
	Message    string
	IsResolved bool
	CreatedAt  time.Time  // gorm.DeletedAtではなく通常のtime.Time
	ResolvedAt *time.Time // ポインタでnullable
}

type AlertConfigModel struct {
	ID          string `gorm:"primaryKey"`
	StockID     string
	MinQuantity int
	MaxQuantity int
	IsActive    bool
	UpdatedAt   time.Time
}

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	// dsnを明示的に定義
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

	// 開発環境のみ: テーブルを削除して再作成
	sqlDB, _ := db.DB()
	sqlDB.Exec("DROP TABLE IF EXISTS alert_models CASCADE")
	sqlDB.Exec("DROP TABLE IF EXISTS alert_config_models CASCADE")

	// マイグレーション実行
	if err := db.AutoMigrate(&AlertModel{}, &AlertConfigModel{}); err != nil {
		log.Printf("Migration error: %v", err)
		return nil, err
	}

	// テーブル名を確認するためのデバッグログ
	var tables []string
	db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables)
	log.Printf("Available tables: %v", tables)

	return db, nil
}
