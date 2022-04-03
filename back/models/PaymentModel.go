package models

import (
	"errors"
	_ "github.com/lib/pq"
)

type Payment struct {
	Transaction   string `json:"transaction"`
	Request_id    string `json:"request_id"`
	Currency      string `json:"currency"`
	Provider      string `json:"provider"`
	Amount        int    `json:"amount"`
	Payment_dt    int    `json:"payment_dt"`
	Bank          string `json:"bank"`
	Delivery_cost int    `json:"delivery_cost"`
	Goods_total   int    `json:"goods_total"`
	Custom_fee    int    `json:"custom_fee"`
}

func (p *Payment) MissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case p.Transaction == "":
		return
	case p.Request_id == "":
		return
	case p.Currency == "":
		return
	case p.Provider == "":
		return
	case p.Amount == 0:
		return
	case p.Payment_dt == 0:
		return
	case p.Bank == "":
		return
	case p.Delivery_cost == 0:
		return
	case p.Goods_total == 0:
		return
	case p.Custom_fee == 0:
		return
	default:
		return nil
	}
}
