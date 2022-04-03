package models

import (
	"errors"
	_ "github.com/lib/pq"
)

type Delivery struct {
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Zip     *string `json:"zip"`
	City    *string `json:"city"`
	Address *string `json:"address"`
	Region  *string `json:"region"`
	Email   *string `json:"email"`
}

func (d *Delivery) MissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case d.Name == nil:
		return
	case d.Phone == nil:
		return
	case d.Zip == nil:
		return
	case d.City == nil:
		return
	case d.Address == nil:
		return
	case d.Region == nil:
		return
	case d.Email == nil:
		return
	default:
		return nil
	}
}
