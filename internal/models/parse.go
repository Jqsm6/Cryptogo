package models

type ETHBalance struct {
	ETH struct {
		Balance int `json:"balance"`
	} `json:"ETH"`
}