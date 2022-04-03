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

func Publisher(nats *nats.Connection) {
	files, err := ioutil.ReadDir("back/test_json/examples")
	if err != nil {
		return
	}
	for _, file := range files {
		name := fmt.Sprintf("back/test_json/examples/%s", file.Name())
		bytes, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		nats.Publish(bytes)
	}
}
func main() {
	db, err := db.OpenConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cache := cache2.CreateCache()
	db.Recovery(cache)

	natsConn, err := nats.NewConn()

	if err != nil {
		log.Fatal(err)
	}

	defer natsConn.NatsConnection.Close()

	output := make(chan []byte)

	subscription, err := natsConn.NewSub(output)

	if err != nil {
		log.Fatal(err)
	}

	defer subscription.Close()

	Publisher(&natsConn)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		var model models.Model
		for bytes := range output {
			model, err = models.NewModel(bytes)

			if err != nil {
				log.Print("Cannot create new model")
			}

			err = model.MissingFields()
			if err != nil {
			} else if _, flag := cache.GetModelCache(*model.Order_uid); flag {
				log.Printf("Model in cache already")
			} else {
				wg.Add(2)
				go db.InsertModel(&wg, model)
				go func() {
					cache.AddModelCache(model)
					wg.Done()
				}()
			}
		}
		wg.Done()
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(output)
		wg.Wait()
		os.Exit(1)
	}()

	serv := model.CreateServer(cache)
	log.Println(http.ListenAndServe(":3000", serv))
}
