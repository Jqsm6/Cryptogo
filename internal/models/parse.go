package models

type ETHBalance struct {
	ETH struct {
		Balance int `json:"balance"`
	} `json:"ETH"`
}

type CountID struct {
	Count int `pg:"count"`
}

type Transaction struct {
	From    string  `json:"from"`
	Value   float64 `json:"value"`
	Success bool    `json:"success"`
}
