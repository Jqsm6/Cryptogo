package models

type Status struct {
	API struct {
		Status struct {
			ETH string `json:"eth"`
		} `json:"Status"`
	} `json:"API"`
}
