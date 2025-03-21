package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"service-2/config"
	"service-2/internal/handler"
	"service-2/internal/middleware"
	"service-2/internal/repository"
	"service-2/internal/service"
)

func main() {
	// 設定を読み込む
	cfg := config.NewConfig()

	// データベース接続
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// リポジトリの初期化
	transactionRepo := repository.NewTransactionRepository(db)
	stockRepo := repository.NewStockRepository(cfg.StockServiceURL)

	// サービスの初期化
	transactionService := service.NewTransactionService(transactionRepo, stockRepo)

	// ハンドラーの初期化
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// ルーターの設定
	router := mux.NewRouter()

	// ミドルウェアの設定
	router.Use(middleware.JSONMiddleware)

	// APIルートの登録
	router.HandleFunc("/api/service-2/transactions", transactionHandler.CreateTransaction).Methods("POST")
	router.HandleFunc("/api/service-2/transactions", transactionHandler.ListTransactions).Methods("GET")
	router.HandleFunc("/api/service-2/transactions/{id}", transactionHandler.GetTransaction).Methods("GET")
	router.HandleFunc("/api/service-2/stocks/{stockId}/transactions", transactionHandler.GetStockTransactions).Methods("GET")
	router.HandleFunc("/api/service-2/stocks/{stockId}/summary", transactionHandler.GetStockSummary).Methods("GET")

	// サーバー起動
	log.Println("Starting transaction service on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.PostgresDSN()

	// 設定をログに出力
	log.Printf("POSTGRES_HOST=%s\n", cfg.PostgresHost)
	log.Printf("POSTGRES_USER=%s\n", cfg.PostgresUser)
	log.Printf("POSTGRES_PASSWORD=%s\n", cfg.PostgresPassword)
	log.Printf("POSTGRES_DB=%s\n", cfg.PostgresDB)

	// データベース接続の試行（リトライ付き）
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database. Retrying in 5 seconds... (Attempt %d/5)", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	// UUID拡張の有効化
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// マイグレーション
	err = repository.Migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
