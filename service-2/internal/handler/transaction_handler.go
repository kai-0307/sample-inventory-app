package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"service-2/internal/domain/model"
	"service-2/internal/service"
)

// トランザクションハンドラー
type TransactionHandler struct {
	transactionService service.ITransactionService
}

// 新しいトランザクションハンドラーを作成
func NewTransactionHandler(transactionService service.ITransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// 新しいトランザクションを作成
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var request struct {
		StockID  string `json:"stock_id"`
		Type     string `json:"type"`
		Quantity int    `json:"quantity"`
		Note     string `json:"note"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionService.CreateTransaction(
		request.StockID,
		request.Type,
		request.Quantity,
		request.Note,
	)

	if err != nil {
		switch err {
		case model.ErrEmptyStockID, model.ErrInvalidQuantity, model.ErrInvalidTransactionType:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// トランザクション一覧を取得
func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit := 100
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	transactions, err := h.transactionService.ListTransactions(limit, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// IDによってトランザクションを取得
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	transaction, err := h.transactionService.GetTransaction(id)
	if err != nil {
		switch err {
		case model.ErrTransactionNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// 在庫IDによってトランザクションを取得
func (h *TransactionHandler) GetStockTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockID := vars["stockId"]

	transactions, err := h.transactionService.GetStockTransactions(stockID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// 在庫の集計情報を取得
func (h *TransactionHandler) GetStockSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockID := vars["stockId"]

	summary, err := h.transactionService.GetStockSummary(stockID)
	if err != nil {
		switch err {
		case model.ErrStockNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
