package models

type PaymentRequest struct {
	Currency    string  `json:"currency"`
	Amount      float64 `json:"amount"`
	ToAddress   string  `json:"to_address"`
	FromAddress string  `json:"from_address"`
}

type PaymentResponse struct {
	ID        string `json:"id"`
	ToAddress string `json:"to_address"`
}

type PaymentInfoRequest struct {
	ID string `json:"id"`
}

type PaymentInfoResponse struct {
	ID          string `json:"id"`
	State       string `json:"state"`
	Currency    string `json:"currency"`
	Amount      string `json:"amount"`
	ToAddress   string `json:"to_address"`
	FromAddress string `json:"from_address"`
}

type PaymentDB struct {
	ID          string `pg:"invoice_id"`
	State       string `pg:"state"`
	Currency    string `pg:"currency"`
	Amount      string `pg:"amount"`
	ToAddress   string `pg:"to_address"`
	FromAddress string `pg:"from_adress"`
}

type PaymentConfirmRequest struct {
	ID     string `json:"id"`
	TxHash string `json:"tx_hash"`
}

type PaymentConfirmResponse struct {
	State string `json:"state"`
}
