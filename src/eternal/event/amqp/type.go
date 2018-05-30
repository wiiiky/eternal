package amqp

import (
	"github.com/streadway/amqp"
)

const (
	METERNAL_TYPE = "eternal/json"
	CONTENT_TYPE  = "application/json"
)

type Publisher struct {
	Conn     *amqp.Connection
	Chan     *amqp.Channel
	Exchange string
	RouteKey string
}

type Subscriber struct {
	Conn     *amqp.Connection
	Chan     *amqp.Channel
	Exchange string
	Queue    string
	RouteKey string
	Consumer string
}

type EventHandler func(string, []byte) bool
