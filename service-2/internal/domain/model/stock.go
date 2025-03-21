package model

// Stock は在庫情報を表す
type Stock struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}
