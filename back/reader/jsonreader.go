package reader

import (
	"L0/back/models"
	"encoding/json"
	"io/ioutil"
)

func JSON_Read(jsonName string) (models.Model, error) {
	model := new(models.Model)

	jsonFile, err := ioutil.ReadFile(jsonName)
	if err != nil {
		return models.Model{}, err
	}
	err = json.Unmarshal(jsonFile, model)
	if err != nil {
		return models.Model{}, err
	}

	err = model.MissingFields()
	if err != nil {
		return models.Model{}, err
	}
	return *model, err
}
