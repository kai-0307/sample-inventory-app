package main

import (
	"log"

	"service-3/config"
	database "service-3/infrastructure/db"
	"service-3/infrastructure/server"
	"service-3/infrastructure/services"
	"service-3/interfaces/handlers"
	"service-3/usecases"
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
}
