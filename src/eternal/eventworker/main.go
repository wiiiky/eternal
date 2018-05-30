package main

import (
	"eternal/config"
	"eternal/event"
	"eternal/logging"
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	config.Init()
	initLogging()
	initDatabase()
	initEvent()
	event.Register(event.KeyAnswerUpvote, handleAnswerUpvote)
	event.Register(event.KeyAnswerDownvote, handleAnswerDownvote)

	ch := make(chan os.Signal)
	<-ch
}

func initEvent() {
	viper.SetDefault("event.amqp.consumer", "eternal")
	amqpURL := viper.GetString("event.amqp.url")
	amqpExchange := viper.GetString("event.amqp.exchange")
	amqpRouteKey := viper.GetString("event.amqp.route_key")
	amqpQueue := viper.GetString("event.amqp.queue")
	amqpConsumer := viper.GetString("event.amqp.consumer")
	event.InitSub(amqpURL, amqpExchange, amqpQueue, amqpRouteKey, amqpConsumer)
}

func initLogging() {
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.output", "stdout")

	logging.Init(viper.GetString("log.format"), viper.GetString("log.level"), viper.GetString("log.output"))
}

func initDatabase() {
	dbURL := viper.GetString("database.url")
	if err := db.Init(dbURL); err != nil {
		log.Fatal("Connecting database failed:", err)
	}
}
