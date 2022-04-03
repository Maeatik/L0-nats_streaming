package models

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"time"
)

type Model struct {
	Order_uid 			*string 	`json:"order_uid"`
	Track_number 		string 		`json:"track_number"`
	Entry    			string   	`json:"entry"`
	Delivery 			Delivery 	`json:"delivery"`
	Payment  			Payment  	`json:"payment"`
	Items    			[]Items  	`json:"items"`
	Locale   			string   	`json:"locale"`
	Internal_signature 	string 		`json:"internal_signature"`
	Customer_id 		string 		`json:"customer_id"`
	Delivery_service 	string 		`json:"delivery_service"`
	Shardkey 			string 		`json:"shardkey"`
	Sm_id 				int 		`json:"sm_id"`
	Date_created 		time.Time 	`json:"date_created"`
	Oof_shard 			string 		`json:"oof_shard"`
}

func NewModel(byteModel []byte) (Model, error) {
	res := new(Model)
	err := json.Unmarshal(byteModel, res)
	return *res, err
}

func (m *Model) Marshal() ([]byte, error) {
	res, err := json.Marshal(m)
	return res, err
}
