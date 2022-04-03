package models

import (
	"errors"
	_ "github.com/lib/pq"
)

type Payment struct {
	Transaction   *string `json:"transaction"`
	Request_id    *string `json:"request_id"`
	Currency      *string `json:"currency"`
	Provider      *string `json:"provider"`
	Amount        *int    `json:"amount"`
	Payment_dt    *int    `json:"payment_dt"`
	Bank          *string `json:"bank"`
	Delivery_cost *int    `json:"delivery_cost"`
	Goods_total   *int    `json:"goods_total"`
	Custom_fee    *int    `json:"custom_fee"`
}

func (p *Payment) MissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case p.Transaction == nil:
		return
	case p.Request_id == nil:
		return
	case p.Currency == nil:
		return
	case p.Provider == nil:
		return
	case p.Amount == nil:
		return
	case p.Payment_dt == nil:
		return
	case p.Bank == nil:
		return
	case p.Delivery_cost == nil:
		return
	case p.Goods_total == nil:
		return
	case p.Custom_fee == nil:
		return
	default:
		return nil
	}
}
