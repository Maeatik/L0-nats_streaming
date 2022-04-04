package test

import (
	"L0/back/db"
	"L0/back/reader"
	"log"
	"sync"
	"testing"
)

//Тесты для работы с json фалами

//Тест на обработку json файла и представление информации в виде набора символов
func TestMarshal(t *testing.T) {

	log.Println("Get test data")
	//Чтение указанного json файла
	jsonEx, err := reader.JSON_Read("./model.json")
	//Если файл не нашелся - тест провален
	if err != nil {
		t.FailNow()
	}
	//Чтение стурктуры Model и превращение ее в набор символов
	resByte, _ := jsonEx.Marshal()
	t.Log(string(resByte))
}

//Тест на добавление данных в БД из json файла
func TestInsertModel(t *testing.T) {
	//Объявление ожидания
	var wg sync.WaitGroup
	//Увеличение счетчика ожидания на 1
	wg.Add(1)
	//Установление соединения с БД
	db, err := db.OpenConnection()
	//Если соединение уставновить не удалось - тест провален
	if err != nil {
		t.FailNow()
	}
	//Отложенное закрытие БД
	defer db.Db.Close()

	//Чтение данных из json-файла
	jsonEx, err := reader.JSON_Read("./examples/model5.json")
	//Если файл не нашелся - тест провален
	if err != nil {
		t.FailNow()
	}
	//Добвление считанных данных в БД
	err = db.InsertModel(&wg, jsonEx)
	if err != nil {
		//Если данные не добавились в БД - тест провален
		t.FailNow()
	}
	//Ожидание пока счетчик не обнулиться
	wg.Wait()
}
