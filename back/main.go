package main

import (
	cache2 "L0/back/cache"
	"L0/back/db"
	"L0/back/model-service"
	"L0/back/models"
	"L0/back/nats"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Publisher - Публикатор. Открывает директорию и читает находящиеся в ней файлы
func Publisher(nats *nats.Connection) {
	//Открывает директорию back/test_json/examples и читает находящиеся в ней файлы
	files, err := ioutil.ReadDir("back/test_json/examples")
	if err != nil {
		return
	}
	//Обрабатывает все найденные документы по очереди
	for _, file := range files {
		name := fmt.Sprintf("back/test_json/examples/%s", file.Name())
		bytes, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		//Отправляет на публикацию подписчику
		nats.Publish(bytes)
	}
}

//Главная функция main
func main() {
	//Установка соединения с БД
	db, err := db.OpenConnection()
	if err != nil {
		log.Fatal(err)
	}
	//Отложенное закрытие соединения с БД
	defer db.Close()

	//Создание кеша
	cache := cache2.CreateCache()
	//Добавления данных в кеш
	db.Recovery(cache)

	//Установка соединения с nats-streaming
	natsConn, err := nats.NewConn()

	if err != nil {
		log.Fatal(err)
	}
	//Отложенное закрытие соединения с nats-streaming
	defer natsConn.NatsConnection.Close()
	//создание канала подписчика для получения данных
	output := make(chan []byte)
	//создание подписчика с указанием для него его канала
	subscription, err := natsConn.NewSub(output)

	if err != nil {
		log.Fatal(err)
	}
	//Отложенное закрытие подпичика
	defer subscription.Close()
	//Создание публикатора, с ссылкой на установленное соединение
	Publisher(&natsConn)

	//Объявление ожидания
	var wg sync.WaitGroup
	//Увелечение считчика ожидания на 1
	wg.Add(1)
	//Запуск горутины для работы с кешем
	go func() {
		//Объявление модели
		var model models.Model
		//Обработка символов в канале подписчика
		for bytes := range output {
			//Создание новой модели по данным, полученных из канала подписчика
			model, err = models.NewModel(bytes)

			if err != nil {
				log.Print("Cannot create new model")
			}
			//Проверка на целостность json файлов
			err = model.MissingFields()
			if err != nil {
			} else if _, flag := cache.GetModelCache(*model.Order_uid); flag {
				//Если ошибки в файлах нет, выдаются данные по заданному order_uid
				//Если модель уже в кеше, выходит уведомление об этом
				log.Printf("Model in cache already")
			} else {
				//Если модель не в кеше, счетчик ожидания увеличивается на 2
				wg.Add(2)
				//Данные добавляются в таблицы, с указанием на модель и указанныый счетчик
				go db.InsertModel(&wg, model)
				go func() {
					//В кеш добавляется новая модель
					cache.AddModelCache(model)
					//Счетчик ожидания уменьшается на 1
					wg.Done()
				}()
			}
		}
		//Счетчик уменьшается на 1
		wg.Done()
	}()
	//Завершение программы
	//Создание канала для обработки сиганолов
	c := make(chan os.Signal)
	//Сигнал оповещающий о завершении программы
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		//Получение данных из канала
		<-c
		//Закрытие канала подписчика
		close(output)
		//Ожидание остановки горутин
		wg.Wait()
		//Немедленное остановление программы
		os.Exit(1)
	}()
	//Получение данных для сервера
	serv := model.CreateServer(cache)
	//Запуск бека на порте :3000
	log.Println(http.ListenAndServe(":3000", serv))
}
