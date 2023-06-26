package dto

type Response struct {
	Success    bool   `json:"success"`
	Cause      string `json:"cause"`
	Id         int    `json:"id"`
	ExternalID string `json:"externalID"`
	Amount     int    `json:"amount"`
	CardNumber string `json:"cardNumber"`
}
