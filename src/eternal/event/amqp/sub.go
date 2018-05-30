package amqp

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strings"
)

var handlers = map[string]EventHandler{}
var subscriber *Subscriber = nil

func Register(routeKey string, handler EventHandler) {
	if subscriber == nil {
		log.Error("InitSub() must be called before Register()")
		return
	}
	key := strings.Replace(subscriber.RouteKey, "#", routeKey, -1)
	handlers[key] = handler
}

func InitSub(url, exchange, queue, routeKey, consumer string) {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(exchange, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(queue, routeKey, exchange, false, nil)
	if err != nil {
		panic(err)
	}

	ch.Qos(2, 0, false)

	subscriber = &Subscriber{
		Conn:     conn,
		Chan:     ch,
		Exchange: exchange,
		Queue:    queue,
		RouteKey: routeKey,
		Consumer: consumer,
	}

	go func() {
		notifies, err := ch.Consume(queue, consumer, false, false, true, false, nil)
		if err != nil {
			panic(err)
		}

		for notify := range notifies {
			routeKey := notify.RoutingKey
			body := notify.Body
			handleFunc := handlers[routeKey]
			log.Infof("AMQP Receive %s", routeKey)
			if handleFunc == nil {
				notify.Ack(false)
			} else {
				notify.Ack(handleFunc(routeKey, body))
			}
		}
	}()
}

func (s *Subscriber) Run() {
	notifies, err := s.Chan.Consume(s.Queue, s.Consumer, false, false, true, false, nil)
	if err != nil {
		panic(err)
	}

	for notify := range notifies {
		routeKey := notify.RoutingKey
		body := notify.Body
		handler := handlers[routeKey]
		log.Infof("AMQP Receive %s", routeKey)
		if handler == nil {
			log.Warnf("AMQP Unknown Key %s", routeKey)
			notify.Ack(false)
		} else {
			notify.Ack(handler(routeKey, body))
		}
	}
}
