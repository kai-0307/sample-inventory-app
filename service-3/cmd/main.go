package main

import (
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

type Alert struct {
    ID           string    `json:"id" gorm:"primaryKey"`
    StockID      string    `json:"stock_id"`
    Type         string    `json:"type"`      // "low_stock", "excess_stock", etc.
    Message      string    `json:"message"`
    IsResolved   bool      `json:"is_resolved"`
    CreatedAt    time.Time `json:"created_at"`
    ResolvedAt   *time.Time `json:"resolved_at"`
}

type AlertConfig struct {
    ID           string    `json:"id" gorm:"primaryKey"`
    StockID      string    `json:"stock_id"`
    MinQuantity  int       `json:"min_quantity"`
    MaxQuantity  int       `json:"max_quantity"`
    IsActive     bool      `json:"is_active"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type StockReport struct {
    StockID      string    `json:"stock_id"`
    CurrentStock int       `json:"current_stock"`
    MinQuantity  int       `json:"min_quantity"`
    MaxQuantity  int       `json:"max_quantity"`
    AlertCount   int       `json:"alert_count"`
    LastAlert    time.Time `json:"last_alert"`
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
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
		)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // マイグレーション実行
    db.AutoMigrate(&Alert{}, &AlertConfig{})

    server := &Server{
        db:     db,
        router: mux.NewRouter(),
    }

    server.routes()
    return server, nil
}

func (s *Server) routes() {
    // アラート関連
    s.router.HandleFunc("/api/service-3/alerts", s.listAlerts).Methods("GET")
    s.router.HandleFunc("/api/service-3/alerts/{id}", s.getAlert).Methods("GET")
    s.router.HandleFunc("/api/service-3/alerts/{id}/resolve", s.resolveAlert).Methods("POST")

    // アラート設定関連
    s.router.HandleFunc("/api/service-3/configs", s.createAlertConfig).Methods("POST")
    s.router.HandleFunc("/api/service-3/configs/{stockId}", s.getAlertConfig).Methods("GET")
    s.router.HandleFunc("/api/service-3/configs/{stockId}", s.updateAlertConfig).Methods("PUT")

    // レポート関連
    s.router.HandleFunc("/api/service-3/reports/stocks", s.generateStockReport).Methods("GET")
    s.router.HandleFunc("/api/service-3/reports/alerts", s.generateAlertReport).Methods("GET")
}

func (s *Server) getAlert(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var alert Alert
    if err := s.db.First(&alert, "id = ?", id).Error; err != nil {
        http.Error(w, "alert not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(alert)
}

func (s *Server) getAlertConfig(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    stockID := vars["stockId"]

    var config AlertConfig
    if err := s.db.First(&config, "stock_id = ?", stockID).Error; err != nil {
        http.Error(w, "alert config not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(config)
}

func (s *Server) updateAlertConfig(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    stockID := vars["stockId"]

    var config AlertConfig
    if err := s.db.First(&config, "stock_id = ?", stockID).Error; err != nil {
        http.Error(w, "alert config not found", http.StatusNotFound)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    config.UpdatedAt = time.Now()

    if err := s.db.Save(&config).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(config)
}

func (s *Server) createAlertConfig(w http.ResponseWriter, r *http.Request) {
    var config AlertConfig
    if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    config.UpdatedAt = time.Now()
    config.IsActive = true

    if result := s.db.Create(&config); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(config)
}

func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
    var alerts []Alert

    query := r.URL.Query()
    resolved := query.Get("resolved")
    stockID := query.Get("stock_id")

    db := s.db
    if resolved != "" {
        isResolved := resolved == "true"
        db = db.Where("is_resolved = ?", isResolved)
    }

    if stockID != "" {
        db = db.Where("stock_id = ?", stockID)
    }

    if result := db.Order("created_at desc").Find(&alerts); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(alerts)
}

func (s *Server) resolveAlert(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var alert Alert
    if err := s.db.First(&alert, "id = ?", id).Error; err != nil {
        http.Error(w, "alert not found", http.StatusNotFound)
        return
    }

    now := time.Now()
    alert.IsResolved = true
    alert.ResolvedAt = &now

    if err := s.db.Save(&alert).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(alert)
}

func (s *Server) generateStockReport(w http.ResponseWriter, r *http.Request) {
    var reports []StockReport

    rows, err := s.db.Raw(`
        SELECT
            ac.stock_id,
            ac.min_quantity,
            ac.max_quantity,
            COUNT(a.id) as alert_count,
            MAX(a.created_at) as last_alert
        FROM
            alert_configs ac
            LEFT JOIN alerts a ON ac.stock_id = a.stock_id
        WHERE
            ac.is_active = true
        GROUP BY
            ac.stock_id, ac.min_quantity, ac.max_quantity
    `).Rows()

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var report StockReport
        rows.Scan(
            &report.StockID,
            &report.MinQuantity,
            &report.MaxQuantity,
            &report.AlertCount,
            &report.LastAlert,
        )
        reports = append(reports, report)
    }

    // 現在の在庫数を取得（Stock Serviceから）
    for i := range reports {
        reports[i].CurrentStock = s.fetchCurrentStock(reports[i].StockID)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reports)
}

func (s *Server) generateAlertReport(w http.ResponseWriter, r *http.Request) {
    type AlertSummary struct {
        TotalAlerts     int64     `json:"total_alerts"`
        ResolvedAlerts  int64     `json:"resolved_alerts"`
        ActiveAlerts    int64     `json:"active_alerts"`
        LastAlertDate   time.Time `json:"last_alert_date"`
    }

    var summary AlertSummary

    // 合計アラート数
    s.db.Model(&Alert{}).Count(&summary.TotalAlerts)

    // 解決済みアラート数
    s.db.Model(&Alert{}).Where("is_resolved = ?", true).Count(&summary.ResolvedAlerts)

    // アクティブなアラート数
    s.db.Model(&Alert{}).Where("is_resolved = ?", false).Count(&summary.ActiveAlerts)

    // 最新のアラート日時
    var lastAlert Alert
    s.db.Order("created_at desc").First(&lastAlert)
    summary.LastAlertDate = lastAlert.CreatedAt

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(summary)
}

// Stock Serviceから現在の在庫数を取得
func (s *Server) fetchCurrentStock(stockID string) int {
    // 実際の実装ではHTTPリクエストでStock Serviceから取得
    // ここではダミー値を返す
    return 100
}

func main() {
    server, err := NewServer()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Starting alert service on :8080")
    if err := http.ListenAndServe(":8080", server.router); err != nil {
        log.Fatal(err)
    }
}
