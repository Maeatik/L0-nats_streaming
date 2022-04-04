package nats

import (
	"github.com/nats-io/stan.go"
	"log"
)

//Константы для работы с nats-streaming
const (
	nats_cluster = "test-cluster"
	nats_subj    = "models"
	nats_client  = "client"
)

// Connection - Структура Соединения с nats-streaming
type Connection struct {
	NatsConnection stan.Conn
}

// NewConn - Новое соединение
func NewConn() (Connection, error) {
	//Создается новый объект структуры Connection
	natsConn := new(Connection)
	var err error
	//Устанавливается соединение с nats-streaming, кластер указан там же
	natsConn.NatsConnection, err = stan.Connect(nats_cluster, nats_client)

	return *natsConn, err
}

// NewSub - Создание нового подписчика и канала для него
func (n *Connection) NewSub(output chan<- []byte) (stan.Subscription, error) {
	//Оформление подписка для объекта
	sub, err := n.NatsConnection.Subscribe(nats_subj, func(msg *stan.Msg) {
		//Установление канала для подписчика
		output <- msg.Data
	},
		//определение метода доставки всех доступных сообщений для подпичика
		stan.DeliverAllAvailable())
	return sub, err
}

// Publish - Публикация данных из json файлов, которые представлены в виде наборов символов
func (n Connection) Publish(jsonFile []byte) {
	//Отправка сообщения подписчику
	err := n.NatsConnection.Publish(nats_subj, jsonFile)
	if err != nil {
		log.Fatal(err)
	}
}
