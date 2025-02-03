package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"service-1/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() (*gorm.DB, error) {
  // 環境変数の取得とデフォルト値の設定
  host := os.Getenv("POSTGRES_HOST")
  user := os.Getenv("POSTGRES_USER")
  password := os.Getenv("POSTGRES_PASSWORD")
  dbname := os.Getenv("POSTGRES_DB")

  // デフォルト値の設定
  if host == "" {
      host = "postgres"
  }
  if user == "" {
      user = "stockapp"
  }
  if password == "" {
      password = "stockapp"
  }
  if dbname == "" {
      dbname = "stockapp"
  }

  dsn := fmt.Sprintf(
      "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tokyo",
      host,
      user,
      password,
      dbname,
  )

  // デバッグ用にDSN情報を出力
  log.Printf("Connecting to database with DSN: %s", dsn)

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

  // マイグレーション
  if err := db.AutoMigrate(&models.Stock{}); err != nil {
      return nil, fmt.Errorf("failed to migrate database: %v", err)
  }

  return db, nil
}
