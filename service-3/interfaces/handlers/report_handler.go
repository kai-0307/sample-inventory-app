package handlers

import (
	"encoding/json"
	"net/http"

	"service-3/usecases"
)

type ReportHandler struct {
	useCase *usecases.ReportUseCase
}

func NewReportHandler(useCase *usecases.ReportUseCase) *ReportHandler {
	return &ReportHandler{useCase: useCase}
}

func (h *ReportHandler) GenerateStockReport(w http.ResponseWriter, r *http.Request) {
	reports, err := h.useCase.GenerateStockReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func (h *ReportHandler) GenerateAlertReport(w http.ResponseWriter, r *http.Request) {
	summary, err := h.useCase.GenerateAlertReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
