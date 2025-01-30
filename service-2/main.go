// service-2

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StockID     string    `json:"stock_id"`
	Type        string    `json:"type"`      // "in" or "out"
	Quantity    int       `json:"quantity"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
}

type TransactionSummary struct {
    StockID     string `json:"stock_id"`
    TotalIn     int    `json:"total_in"`
    TotalOut    int    `json:"total_out"`
    Balance     int    `json:"balance"`
}

type Server struct {
    db     *gorm.DB
    router *mux.Router
}

func NewServer() (*Server, error) {
	// データベース接続設定
	dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tokyo",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_PASSWORD"),
	)

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
			return nil, fmt.Errorf("failed to connect to database after 5 attempts: %v", err)
	}

	// UUID拡張の有効化
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// マイグレーション実行
	if err := db.AutoMigrate(&Transaction{}); err != nil {
			return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	server := &Server{
			db:     db,
			router: mux.NewRouter(),
	}

	server.routes()
	return server, nil
}

func (s *Server) routes() {
    s.router.HandleFunc("/api/service-2/transactions", s.createTransaction).Methods("POST")
    s.router.HandleFunc("/api/service-2/transactions", s.listTransactions).Methods("GET")
    s.router.HandleFunc("/api/service-2/transactions/{id}", s.getTransaction).Methods("GET")
    s.router.HandleFunc("/api/service-2/stocks/{stockId}/transactions", s.getStockTransactions).Methods("GET")
    s.router.HandleFunc("/api/service-2/stocks/{stockId}/summary", s.getStockSummary).Methods("GET")
}

func (s *Server) createTransaction(w http.ResponseWriter, r *http.Request) {
	// トランザクションの開始
	tx := s.db.Begin()
	if tx.Error != nil {
			http.Error(w, tx.Error.Error(), http.StatusInternalServerError)
			return
	}

	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
	}

	// バリデーション
	if transaction.StockID == "" || transaction.Quantity == 0 {
			tx.Rollback()
			http.Error(w, "invalid transaction data", http.StatusBadRequest)
			return
	}

	// トランザクションタイプの検証
	if transaction.Type != "in" && transaction.Type != "out" {
			tx.Rollback()
			http.Error(w, "invalid transaction type", http.StatusBadRequest)
			return
	}

	// 出庫の場合は負の値に変換
	quantity := transaction.Quantity
	if transaction.Type == "out" {
			quantity = -quantity
	}

	transaction.CreatedAt = time.Now()

	// データベースに保存
	if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	// Stock Serviceの在庫数を更新
	if err := s.updateStockQuantity(transaction.StockID, quantity); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	// トランザクションのコミット
	if err := tx.Commit().Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func (s *Server) listTransactions(w http.ResponseWriter, r *http.Request) {
    var transactions []Transaction

    // クエリパラメータの取得
    query := r.URL.Query()
    limit := 100 // デフォルトの取得件数

    // 日付範囲フィルター
    startDate := query.Get("start_date")
    endDate := query.Get("end_date")

    db := s.db
    if startDate != "" && endDate != "" {
        db = db.Where("created_at BETWEEN ? AND ?", startDate, endDate)
    }

    result := db.Order("created_at desc").Limit(limit).Find(&transactions)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

func (s *Server) getTransaction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var transaction Transaction
    result := s.db.First(&transaction, "id = ?", id)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transaction)
}

func (s *Server) getStockTransactions(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    stockID := vars["stockId"]

    var transactions []Transaction
    result := s.db.Where("stock_id = ?", stockID).
        Order("created_at desc").
        Find(&transactions)

    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
}

func (s *Server) getStockSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockID := vars["stockId"]

	var summary TransactionSummary
	summary.StockID = stockID

	// Stock Serviceから現在の在庫数を取得
	resp, err := http.Get(fmt.Sprintf("http://service-1:8080/api/service-1/stocks/%s", stockID))
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	defer resp.Body.Close()

	var stock struct {
			Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	// 入庫合計の計算
	s.db.Model(&Transaction{}).
			Where("stock_id = ? AND type = ?", stockID, "in").
			Select("COALESCE(SUM(quantity), 0)").
			Scan(&summary.TotalIn)

	// 出庫合計の計算
	s.db.Model(&Transaction{}).
			Where("stock_id = ? AND type = ?", stockID, "out").
			Select("COALESCE(SUM(quantity), 0)").
			Scan(&summary.TotalOut)

	// 現在の在庫数を設定
	summary.Balance = stock.Quantity

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// 新しい関数を追加
func (s *Server) updateStockQuantity(stockID string, quantity int) error {
	url := fmt.Sprintf("http://service-1:8080/api/service-1/stocks/%s/quantity", stockID)
	payload := map[string]interface{}{
			"adjustment": quantity,
			"note": "Transaction update",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
			return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
			return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
			return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to update stock quantity: %d", resp.StatusCode)
	}

	return nil
}

func main() {
    server, err := NewServer()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Starting transaction service on :8080")
    if err := http.ListenAndServe(":8080", server.router); err != nil {
        log.Fatal(err)
    }
}
