package main

import (
	"eternal/config"
	"eternal/event"
	"eternal/logging"
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
	"os"
)

const APPNAME = "eventworker"

func main() {
	config.Init(APPNAME)
	initLogging()
	initDatabase()
	initEvent()
	event.Register(event.KeyAnswerUpvote, handleAnswerUpvote)
	event.Register(event.KeyAnswerDownvote, handleAnswerDownvote)

	ch := make(chan os.Signal)
	<-ch
}

func initEvent() {
	amqpURL := config.GetString("event.amqp.url")
	amqpExchange := config.GetString("event.amqp.exchange")
	amqpRouteKey := config.GetString("event.amqp.route_key")
	amqpQueue := config.GetString("event.amqp.queue")
	amqpConsumer := config.GetStringDefault("event.amqp.consumer", "eventworker")
	event.InitSub(amqpURL, amqpExchange, amqpQueue, amqpRouteKey, amqpConsumer)
}

func initLogging() {
	format := config.GetStringDefault("log.format", "json")
	level := config.GetStringDefault("log.level", "info")
	output := config.GetStringDefault("log.output", "stdout")

	logging.Init(format, level, output)
}

func initDatabase() {
	dbURL := config.GetString("database.url")
	if dbURL == "" {
		log.Fatal("**CONFIG** database.url not found")
	} else if err := db.Init(dbURL); err != nil {
		log.Fatal("Connecting database failed:", err)
	}
}
