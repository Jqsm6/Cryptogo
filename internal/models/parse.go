package models

type Count struct {
	Count int `pg:"count"`
}

type InfoDB struct {
	Amount      string `pg:"amount"`
	ToAddress   string `pg:"to_address"`
	FromAddress string `pg:"from_adress"`
}