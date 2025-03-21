package model

// TransactionSummary は在庫のトランザクション集計を表す
type TransactionSummary struct {
	StockID  string `json:"stock_id"`
	TotalIn  int    `json:"total_in"`
	TotalOut int    `json:"total_out"`
	Balance  int    `json:"balance"`
}
