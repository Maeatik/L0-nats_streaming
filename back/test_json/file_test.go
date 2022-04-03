package test

import (
	"L0/back/db"
	"L0/back/reader"
	"log"
	"sync"
	"testing"
)

func TestMarshal(t *testing.T) {


	log.Println("Get test data")
	jsonEx, err := reader.JSON_Read("./model.json")
	if err != nil {
		t.FailNow()
	}

	resByte, _ := jsonEx.Marshal()
	t.Log(string(resByte))
}

func TestInsertModel(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)

	db, err := db.OpenConnection()
	if err != nil {
		t.FailNow()
	}
	defer db.Db.Close()

	jsonEx, err := reader.JSON_Read("./examples/model5.json")
	if err != nil {
		t.FailNow()
	}

	err = db.InsertModel(&wg, jsonEx)
	if err != nil {
		t.FailNow()
	}
	wg.Wait()
}
