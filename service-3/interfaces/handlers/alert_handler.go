package handlers

import (
	"encoding/json"
	"net/http"

	"service-3/usecases"

	"github.com/gorilla/mux"
)

type AlertHandler struct {
	useCase *usecases.AlertUseCase
}

func NewAlertHandler(useCase *usecases.AlertUseCase) *AlertHandler {
	return &AlertHandler{useCase: useCase}
}

func (h *AlertHandler) GetAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	alert, err := h.useCase.GetAlert(id)
	if err != nil {
		http.Error(w, "alert not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}

func (h *AlertHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	resolved := query.Get("resolved")
	stockID := query.Get("stock_id")

	var isResolved *bool
	if resolved != "" {
		resolvedBool := resolved == "true"
		isResolved = &resolvedBool
	}

	alerts, err := h.useCase.ListAlerts(isResolved, stockID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

func (h *AlertHandler) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	alert, err := h.useCase.ResolveAlert(id)
	if err != nil {
		http.Error(w, "alert not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}
