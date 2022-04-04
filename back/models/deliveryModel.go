package models

import (
	"errors"
	_ "github.com/lib/pq"
)

// Delivery - структура таблицы delivery. Каждое поле имеет ссылочный тип,
//чтобы проверять на наличие в json файле нужную структуру
type Delivery struct {
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
	Zip     *string `json:"zip"`
	City    *string `json:"city"`
	Address *string `json:"address"`
	Region  *string `json:"region"`
	Email   *string `json:"email"`
}

// MissingFields Проверка на наличие записи того или иного поля
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
	//Если все поля имеют запись и ошибок не возникло - возвращается отсутствие ошибки
	default:
		return nil
	}
}
