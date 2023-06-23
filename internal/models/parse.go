package models

type ETHBalance struct {
	ETH struct {
		Balance int `json:"balance"`
	} `json:"ETH"`
}

type Count struct {
	Count int `pg:"count"`
}

type ETHTransaction struct {
	Hash  string  `json:"hash"`
	From  string  `json:"from"`
	Value float64 `json:"value"`
}
