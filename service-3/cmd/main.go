package main

import (
	"log"
	"time"

	"service-3/config"
	database "service-3/infrastructure/db"
	"service-3/infrastructure/server"
	"service-3/infrastructure/services"
	"service-3/interfaces/handlers"
	"service-3/usecases"

	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize database connection
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	alertRepo := database.NewAlertRepository(db)
	configRepo := database.NewAlertConfigRepository(db)

	// Initialize external services
	stockService := services.NewStockService()

	// Initialize use cases
	alertUseCase := usecases.NewAlertUseCase(alertRepo)
	configUseCase := usecases.NewAlertConfigUseCase(configRepo)
	reportUseCase := usecases.NewReportUseCase(alertRepo, configRepo, stockService)

	// Initialize HTTP handlers
	alertHandler := handlers.NewAlertHandler(alertUseCase)
	configHandler := handlers.NewAlertConfigHandler(configUseCase)
	reportHandler := handlers.NewReportHandler(reportUseCase)

	// Start HTTP server
	httpServer := server.NewHTTPServer(alertHandler, configHandler, reportHandler)

	log.Println("Starting alert service on :8080")
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}

	// サーバー起動前にテストデータ追加する(不要であればコメントアウトしておく)
	if err := insertTestData(db); err != nil {
		log.Printf("Warning: Failed to insert test data: %v", err)
	}
}

func debugDatabase(db *gorm.DB) {
	// テーブル一覧を取得
	var tables []string
	db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables)
	log.Printf("Database tables: %v", tables)

	// AlertConfigModelsのレコード数を確認
	var count int64
	db.Model(&database.AlertConfigModel{}).Count(&count)
	log.Printf("AlertConfigModel count: %d", count)

	// サンプルのレコードを取得して表示
	var configs []database.AlertConfigModel
	db.Limit(5).Find(&configs)
	for i, config := range configs {
		log.Printf("Config %d: ID=%s, StockID=%s", i, config.ID, config.StockID)
	}
}

// テストデータ挿入関数
func insertTestData(db *gorm.DB) error {
	// アラート設定のテストデータ
	configModel := database.AlertConfigModel{
		ID:          "config123",
		StockID:     "stock001",
		MinQuantity: 10,
		MaxQuantity: 100,
		IsActive:    true,
		UpdatedAt:   time.Now(),
	}

	// アラートのテストデータ
	alertModel := database.AlertModel{
		ID:         "alert123",
		StockID:    "stock001",
		Type:       "low_stock",
		Message:    "Stock is below minimum threshold",
		IsResolved: false,
		CreatedAt:  time.Now(),
	}

	// データベースに挿入
	if err := db.Create(&configModel).Error; err != nil {
		return err
	}

	if err := db.Create(&alertModel).Error; err != nil {
		return err
	}

	return nil
}
