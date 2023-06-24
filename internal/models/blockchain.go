package models

type ETHTransaction struct {
	Result struct {
		From  string `json:"from"`
		To    string `json:"to"`
		Value string `json:"value"`
	} `json:"result"`
}
