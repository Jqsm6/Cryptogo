package models

type PaymentRequest struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type PaymentResponse struct {
	ID        string `json:"id"`
	ToAddress string `json:"to_address"`
}

type PaymentInfoRequest struct {
	ID string `json:"id"`
}

type PaymentInfoResponse struct {
	State     string `json:"state"`
	Currency  string `json:"currency"`
	Amount    string `json:"amount"`
	ToAddress string `json:"to_address"`
}

type PaymentDB struct {
	ID         string `pg:"invoice_id"`
	State      string `pg:"state"`
	Currency   string `pg:"currency"`
	Amount     string `pg:"amount"`
	ToAddress  string `pg:"to_address"`
	PrivateKey string `pg:"private_key"`
}
