package models

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
