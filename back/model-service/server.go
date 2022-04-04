package model

import (
	"L0/back/cache"
	"fmt"
	"log"
	"net/http"
)

// Server - Структура сервера
type Server struct {
	http.Handler
}

// CreateServer - Создание сервера
func CreateServer(modelMap *cache.Cache) Server {
	//Создается новый объект структуры
	s := new(Server)
	//Создание роутера
	r := *http.NewServeMux()
	//Объявление запросов
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//Обработка полученного метода с фронта
		switch r.Method {
		//Открытие фронта
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Origin", "*")
		//Обработка введенного на фронте order_uid и получение по нему модели
		case http.MethodGet:
			fmt.Println(1)
			getModel, flag := modelMap.GetModelCache(r.Header.Get("order_uid"))
			fmt.Print(getModel)
			if !flag {
				w.WriteHeader(http.StatusNoContent)
				//Вывод полученной модели на страницу фронта
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Access-Control-Allow-Headers", "*")
				w.Header().Set("Access-Control-Allow-Methods", "*")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				//Получает данные в виде набора символов
				pGetModel, err := getModel.Marshal()
				if err != nil {
					log.Printf("Err")
				}
				w.Write(pGetModel)
			}
		//Если возникла ошибка - вывод ошибки 405
		default:
			http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		}
	})
	//Сервер получает ссылку на роутер
	s.Handler = &r
	return *s
}
