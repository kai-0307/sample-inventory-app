package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"service-1/internal/domain/models"
	"service-1/internal/usecase"

	"github.com/gorilla/mux"
)

type StockHandler struct {
    stockInteractor *usecase.StockInteractor
}

func NewStockHandler(stockInteractor *usecase.StockInteractor) *StockHandler {
    return &StockHandler{
        stockInteractor: stockInteractor,
    }
}

func (h *StockHandler) CreateStock(w http.ResponseWriter, r *http.Request) {
    var stock *models.Stock
    if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.stockInteractor.CreateStock(r.Context(), stock); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(stock)
}

func (h *StockHandler) GetStock(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    stock, err := h.stockInteractor.GetStock(r.Context(), id)
    if err != nil {
        http.Error(w, "stock not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stock)
}

func (h *StockHandler) ListStocks(w http.ResponseWriter, r *http.Request) {
    lowStock := r.URL.Query().Get("low_stock") == "true"

    var stocks []*models.Stock
    var err error

    if lowStock {
        stocks, err = h.stockInteractor.ListLowStocks(r.Context())
    } else {
        stocks, err = h.stockInteractor.ListStocks(r.Context())
    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stocks)
}

func (h *StockHandler) UpdateStockQuantity(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var update struct {
        Adjustment int    `json:"adjustment"`
        Note      string `json:"note"`
    }

    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.stockInteractor.UpdateStockQuantity(r.Context(), id, update.Adjustment); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    stock, err := h.stockInteractor.GetStock(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stock)
}

func (h *StockHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var stock *models.Stock
    stock, err := h.stockInteractor.GetStock(r.Context(), id)
    if err != nil {
        http.Error(w, "stock not found", http.StatusNotFound)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    stock.UpdatedAt = time.Now()

    if err := h.stockInteractor.UpdateStock(r.Context(), stock); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stock)
}

func (h *StockHandler) DeleteStock(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    if err := h.stockInteractor.DeleteStock(r.Context(), id); err != nil {
        http.Error(w, "stock not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
