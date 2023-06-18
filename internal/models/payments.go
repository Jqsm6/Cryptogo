package models

type PaymentRequest struct {
	Currency    string `json:"currency"`
	Amount      string `json:"amount"`
	FromAddress string `json:"from_address"`
}

type PaymentResponse struct {
	ID        string `json:"id"`
	ToAddress string `json:"to_address"`
}

type PaymentInfoRequest struct {
	ID string `json:"id"`
}

type PaymentInfoResponse struct {
	Status    int    `json:"status"`
	Currency  string `json:"currency"`
	Amount    string `json:"amount"`
	ToAddress string `json:"to_address"`
}

type PaymentDB struct {
	ID          string `pg:"id_"`
	Status      int    `pg:"status"`
	Currency    string `pg:"currency"`
	Amount      string `pg:"amount"`
	FromAddress string `pg:"from_address"`
	ToAddress   string `pg:"to_address"`
	PrivateKey  string `pg:"private_key"`
}
