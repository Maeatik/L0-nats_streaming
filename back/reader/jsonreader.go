package reader

import (
	"L0/back/models"
	"encoding/json"
	"io/ioutil"
)

// JSON_Read - Чтение json файлов
func JSON_Read(jsonName string) (models.Model, error) {

	//Создание новой модели
	model := new(models.Model)
	//Открыте json файла из папки back/test-json/examples
	jsonFile, err := ioutil.ReadFile(jsonName)
	if err != nil {
		return models.Model{}, err
	}
	//Читает json файл и парсит его в набор символов
	err = json.Unmarshal(jsonFile, model)
	if err != nil {
		return models.Model{}, err
	}
	//Проверка на пустые поля json-файла
	err = model.MissingFields()
	if err != nil {
		return models.Model{}, err
	}
	//Если все ок, создается новая модель
	return *model, err
}
