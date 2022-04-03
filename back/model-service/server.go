package model

import (
	"L0/back/cache"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	http.Handler
}
func CreateServer(modelMap *cache.Cache) Server{
	s:=new(Server)
	r := *http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Origin", "*")

		case http.MethodGet:
			fmt.Println(1)
			getModel, flag := modelMap.GetModelCache(r.Header.Get("order_uid"))
			fmt.Print(getModel)
			if !flag {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.Header().Set("Content-Type","application/json")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			pGetModel, err := getModel.Marshal()
			if err != nil{
				log.Printf("Err")
			}
			w.Write(pGetModel)
		}
		default:
			http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		}
	})

	s.Handler = &r
	return *s
}

