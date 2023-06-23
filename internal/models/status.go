package models

type Status struct {
	API struct {
		Status struct {
			ETH string `json:"eth"`
			BTC string `json:"btc"`
			BNB string `json:"bnb"`
		} `json:"Status"`
	} `json:"API"`
}
