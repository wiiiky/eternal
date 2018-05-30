package amqp

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strings"
)

var publisher *Publisher = nil

func InitPub(amqpURL, exchange, routeKey string) {
	conn, err := amqp.Dial(amqpURL)
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

	publisher = &Publisher{
		Conn:     conn,
		Chan:     ch,
		Exchange: exchange,
		RouteKey: routeKey,
	}
}

func Pub(key string, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal failed:", err)
		return
	}
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  CONTENT_TYPE,
		Type:         METERNAL_TYPE,
		Body:         body,
	}
	key = strings.Replace(publisher.RouteKey, "#", key, -1)
	err = publisher.Chan.Publish(publisher.Exchange, key, false, false, msg)
	if err != nil {
		log.Error("AMQP Publish failed:", err)
	} else {
		log.Infof("AMQP Publish %s successfully", key)
	}
}
