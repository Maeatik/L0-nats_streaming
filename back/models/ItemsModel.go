package models

import "errors"

type Items struct {
	Chrt_id      int    `json:"chrt_id"`
	Track_number string `json:"track_number"`
	Price        int    `json:"price"`
	Rid          string `json:"rid"`
	Name         string `json:"name"`
	Sale         int    `json:"sale"`
	Size         string `json:"size"`
	Total_price  int    `json:"total_price"`
	Nm_id        int    `json:"nm_id"`
	Brand        string `json:"brand"`
	Status       int    `json:"status"`
}

func (item *Items) MissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case item.Chrt_id == 0:
		return
	case item.Track_number == "":
		return
	case item.Price == 0:
		return
	case item.Rid == "":
		return
	case item.Name == "":
		return
	case item.Sale == 0:
		return
	case item.Size == "":
		return
	case item.Total_price == 0:
		return
	case item.Nm_id == 0:
		return
	case item.Brand == "":
		return
	case item.Status == 0:
		return
	default:
		return nil
	}
}
