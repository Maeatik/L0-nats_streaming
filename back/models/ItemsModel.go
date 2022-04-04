package models

import "errors"

// Items - структура таблицы items. Каждое поле имеет ссылочный тип,
//чтобы проверять на наличие в json файле нужную структуру
type Items struct {
	Chrt_id      *int    `json:"chrt_id"`
	Track_number *string `json:"track_number"`
	Price        *int    `json:"price"`
	Rid          *string `json:"rid"`
	Name         *string `json:"name"`
	Sale         *int    `json:"sale"`
	Size         *string `json:"size"`
	Total_price  *int    `json:"total_price"`
	Nm_id        *int    `json:"nm_id"`
	Brand        *string `json:"brand"`
	Status       *int    `json:"status"`
}

// MissingFields Проверка на наличие записи того или иного поля
func (item *Items) MissingFields() (errField error) {
	errField = errors.New("missing field")

	switch {
	case item.Chrt_id == nil:
		return
	case item.Track_number == nil:
		return
	case item.Price == nil:
		return
	case item.Rid == nil:
		return
	case item.Name == nil:
		return
	case item.Sale == nil:
		return
	case item.Size == nil:
		return
	case item.Total_price == nil:
		return
	case item.Nm_id == nil:
		return
	case item.Brand == nil:
		return
	case item.Status == nil:
		return
	//Если все поля имеют запись и ошибок не возникло - возвращается отсутствие ошибки
	default:
		return nil
	}
}
