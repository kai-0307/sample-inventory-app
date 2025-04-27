package handlers

import (
	"encoding/json"
	"net/http"

	"service-3/domain/entities"
	"service-3/usecases"

	"github.com/gorilla/mux"
)

type AlertConfigHandler struct {
	useCase *usecases.AlertConfigUseCase
}

func NewAlertConfigHandler(useCase *usecases.AlertConfigUseCase) *AlertConfigHandler {
	return &AlertConfigHandler{useCase: useCase}
}

func (h *AlertConfigHandler) GetAlertConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockID := vars["stockId"]

	config, err := h.useCase.GetAlertConfig(stockID)
	if err != nil {
		http.Error(w, "alert config not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (h *AlertConfigHandler) CreateAlertConfig(w http.ResponseWriter, r *http.Request) {
	var config entities.AlertConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.useCase.CreateAlertConfig(&config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

func (h *AlertConfigHandler) UpdateAlertConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockID := vars["stockId"]

	// まず現在の設定を取得
	config, err := h.useCase.GetAlertConfig(stockID)
	if err != nil {
		http.Error(w, "alert config not found", http.StatusNotFound)
		return
	}

	// リクエストボディから新しい設定を読み込む
	if err := json.NewDecoder(r.Body).Decode(config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 設定を更新
	if err := h.useCase.UpdateAlertConfig(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
