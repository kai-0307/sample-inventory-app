package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string
}

func NewConfig() *Config {
	config := &Config{
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresPort:     "5432",
	}

	// 環境変数の値を表示（デバッグ用）
	fmt.Printf("POSTGRES_HOST=%s\n", config.PostgresHost)
	fmt.Printf("POSTGRES_USER=%s\n", config.PostgresUser)
	fmt.Printf("POSTGRES_PASSWORD=%s\n", config.PostgresPassword)
	fmt.Printf("POSTGRES_DB=%s\n", config.PostgresDB)

	return config
}
