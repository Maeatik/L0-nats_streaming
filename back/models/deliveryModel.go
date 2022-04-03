package models

import (
	"errors"
	_ "github.com/lib/pq"
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func (d *Delivery) MissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case d.Name == "":
		return
	case d.Phone == "":
		return
	case d.Zip == "":
		return
	case d.City == "":
		return
	case d.Address == "":
		return
	case d.Region == "":
		return
	case d.Email == "":
		return
	default:
		return nil
	}
}
