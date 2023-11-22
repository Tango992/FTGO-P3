package models

type Order struct {
	Amount   int `json:"amount"`
	Discount int `json:"discount"`
}
