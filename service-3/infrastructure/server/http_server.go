package server

import (
	"net/http"

	"service-3/interfaces/handlers"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	router *mux.Router
}

func NewHTTPServer(
	alertHandler *handlers.AlertHandler,
	configHandler *handlers.AlertConfigHandler,
	reportHandler *handlers.ReportHandler,
) *HTTPServer {
	router := mux.NewRouter()

	// アラート関連
	router.HandleFunc("/api/service-3/alerts", alertHandler.ListAlerts).Methods("GET")
	router.HandleFunc("/api/service-3/alerts/{id}", alertHandler.GetAlert).Methods("GET")
	router.HandleFunc("/api/service-3/alerts/{id}/resolve", alertHandler.ResolveAlert).Methods("POST")

	// アラート設定関連
	router.HandleFunc("/api/service-3/configs", configHandler.CreateAlertConfig).Methods("POST")
	router.HandleFunc("/api/service-3/configs/{stockId}", configHandler.GetAlertConfig).Methods("GET")
	router.HandleFunc("/api/service-3/configs/{stockId}", configHandler.UpdateAlertConfig).Methods("PUT")

	// レポート関連
	router.HandleFunc("/api/service-3/reports/stocks", reportHandler.GenerateStockReport).Methods("GET")
	router.HandleFunc("/api/service-3/reports/alerts", reportHandler.GenerateAlertReport).Methods("GET")

	return &HTTPServer{router: router}
}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
