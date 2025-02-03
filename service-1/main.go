// service-1

package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/gorilla/mux"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// type Stock struct {
//     ID          string    `json:"id" gorm:"primaryKey"`
//     Name        string    `json:"name"`
//     Description string    `json:"description"`
//     Quantity    int       `json:"quantity"`
//     Unit        string    `json:"unit"`
//     MinQuantity int       `json:"min_quantity"`
//     CreatedAt   time.Time `json:"created_at"`
//     UpdatedAt   time.Time `json:"updated_at"`
// }

// type Server struct {
//     db     *gorm.DB
//     router *mux.Router
// }

// func NewServer() (*Server, error) {
//     // データベース接続設定
// 		host := os.Getenv("POSTGRES_HOST")
//     if host == "" {
//         host = "postgres" // デフォルト値
//     }
// 		dsn := fmt.Sprintf(
// 			"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tokyo",
// 			os.Getenv("POSTGRES_HOST"),
// 			os.Getenv("POSTGRES_USER"),
// 			os.Getenv("POSTGRES_PASSWORD"),
// 			os.Getenv("POSTGRES_DB"),
// 		)
//     db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//     if err != nil {
//         return nil, err
//     }

//     // マイグレーション実行
//     db.AutoMigrate(&Stock{})

//     server := &Server{
//         db:     db,
//         router: mux.NewRouter(),
//     }

//     server.routes()
//     return server, nil
// }

// func (s *Server) routes() {
//     s.router.HandleFunc("/api/service-1/stocks", s.listStocks).Methods("GET")
//     s.router.HandleFunc("/api/service-1/stocks", s.createStock).Methods("POST")
//     s.router.HandleFunc("/api/service-1/stocks/{id}", s.getStock).Methods("GET")
//     s.router.HandleFunc("/api/service-1/stocks/{id}", s.updateStock).Methods("PUT")
//     s.router.HandleFunc("/api/service-1/stocks/{id}", s.deleteStock).Methods("DELETE")
//     s.router.HandleFunc("/api/service-1/stocks/{id}/quantity", s.updateQuantity).Methods("PATCH")
// }

// func (s *Server) listStocks(w http.ResponseWriter, r *http.Request) {
//     var stocks []Stock

//     // クエリパラメータの取得
//     query := r.URL.Query()
//     lowStock := query.Get("low_stock")

//     // データベースクエリの構築
//     db := s.db
//     if lowStock == "true" {
//         db = db.Where("quantity <= min_quantity")
//     }

//     result := db.Find(&stocks)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(stocks)
// }

// func (s *Server) createStock(w http.ResponseWriter, r *http.Request) {
//     var stock Stock
//     if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     stock.CreatedAt = time.Now()
//     stock.UpdatedAt = time.Now()

//     if result := s.db.Create(&stock); result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusCreated)
//     json.NewEncoder(w).Encode(stock)
// }

// func (s *Server) getStock(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id := vars["id"]

//     var stock Stock
//     if err := s.db.First(&stock, "id = ?", id).Error; err != nil {
//         http.Error(w, "stock not found", http.StatusNotFound)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(stock)
// }

// func (s *Server) updateStock(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id := vars["id"]

//     var stock Stock
//     if err := s.db.First(&stock, "id = ?", id).Error; err != nil {
//         http.Error(w, "stock not found", http.StatusNotFound)
//         return
//     }

//     if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     stock.UpdatedAt = time.Now()

//     if err := s.db.Save(&stock).Error; err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(stock)
// }

// func (s *Server) deleteStock(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id := vars["id"]

//     if err := s.db.Delete(&Stock{}, "id = ?", id).Error; err != nil {
//         http.Error(w, "stock not found", http.StatusNotFound)
//         return
//     }

//     w.WriteHeader(http.StatusNoContent)
// }

// func (s *Server) updateQuantity(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id := vars["id"]

//     type QuantityUpdate struct {
//         Adjustment int    `json:"adjustment"`
//         Note      string `json:"note"`
//     }

//     var update QuantityUpdate
//     if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     var stock Stock
//     tx := s.db.Begin()

//     if err := tx.First(&stock, "id = ?", id).Error; err != nil {
//         tx.Rollback()
//         http.Error(w, "stock not found", http.StatusNotFound)
//         return
//     }

//     // 在庫数の更新
//     newQuantity := stock.Quantity + update.Adjustment
//     if newQuantity < 0 {
//         tx.Rollback()
//         http.Error(w, "insufficient stock", http.StatusBadRequest)
//         return
//     }

//     stock.Quantity = newQuantity
//     stock.UpdatedAt = time.Now()

//     if err := tx.Save(&stock).Error; err != nil {
//         tx.Rollback()
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     // トランザクションのコミット
//     if err := tx.Commit().Error; err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Transaction Serviceへの通知
//     go s.notifyTransactionService(stock.ID, update.Adjustment, update.Note)

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(stock)
// }

// func (s *Server) notifyTransactionService(stockID string, adjustment int, note string) {
//     // Transaction Serviceへの通知処理（実装は省略）
//     log.Printf("Notifying transaction service: stock_id=%s, adjustment=%d, note=%s",
//         stockID, adjustment, note)
// }

// func main() {
//     server, err := NewServer()
//     if err != nil {
//         log.Fatal(err)
//     }

//     log.Println("Starting stock service on :8080")
//     if err := http.ListenAndServe(":8080", server.router); err != nil {
//         log.Fatal(err)
//     }
// }
