package nats

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
)

const  (
	nats_cluster = "test-cluster"
	nats_subj = "models"
	nats_client = "client"
)

type Connection struct {
	NatsConnection stan.Conn
}

func NewConn() (Connection, error) {
	natsConn := new(Connection)
	var err error
	fmt.Print(1)
	natsConn.NatsConnection, err = stan.Connect(nats_cluster, nats_client)
	fmt.Println(2)
	return *natsConn, err
}

func (n *Connection) NewSub(output chan <- []byte)(stan.Subscription, error)  {
	sub, err := n.NatsConnection.Subscribe(nats_subj, func(msg *stan.Msg) {
		output <- msg.Data
	}, stan.DeliverAllAvailable())
	return sub, err
}

func (n Connection) Publish(jsonFile []byte)  {
	err := n.NatsConnection.Publish(nats_subj, jsonFile)
	if err != nil{
		log.Fatal(err)
	}
}
