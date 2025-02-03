package main

import (
	"log"
	"net/http"

	"service-1/internal/infrastructure/database"
	"service-1/internal/interfaces/handlers"
	"service-1/internal/interfaces/repositories"
	"service-1/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {
    // データベース接続
    db, err := database.NewPostgresDB()
    if err != nil {
        log.Fatal(err)
    }

    // リポジトリの初期化
    stockRepo := repositories.NewPostgresStockRepository(db)

    // ユースケースの初期化
    stockInteractor := usecase.NewStockInteractor(stockRepo)

    // ハンドラーの初期化
    stockHandler := handlers.NewStockHandler(stockInteractor)

    // ルーターの設定
    router := mux.NewRouter()

    // サブルーターを作成
    apiRouter := router.PathPrefix("/api/service-1").Subrouter()

    // ルーティングの設定
    apiRouter.HandleFunc("/stocks/{id}", stockHandler.UpdateStock).Methods("PUT")
    apiRouter.HandleFunc("/stocks/{id}", stockHandler.DeleteStock).Methods("DELETE")
    apiRouter.HandleFunc("/stocks", stockHandler.CreateStock).Methods("POST")
    apiRouter.HandleFunc("/stocks", stockHandler.ListStocks).Methods("GET")
    apiRouter.HandleFunc("/stocks/{id}", stockHandler.GetStock).Methods("GET")
    apiRouter.HandleFunc("/stocks/{id}/quantity", stockHandler.UpdateStockQuantity).Methods("PATCH")

    // サーバーの起動
    log.Println("Starting stock service on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal(err)
    }
}
