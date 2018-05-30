package event

import (
	"eternal/event/amqp"
)

func InitPub(amqpURL, exchange, amqpRouteKey string) {
	amqp.InitPub(amqpURL, exchange, amqpRouteKey)
}

func InitSub(url, exchange, queue, routeKey, consumer string) {
	amqp.InitSub(url, exchange, queue, routeKey, consumer)
}

func Publish(key string, data interface{}) {
	go amqp.Pub(key, data)
}

func Register(key string, handler amqp.EventHandler) {
	amqp.Register(key, handler)
}
