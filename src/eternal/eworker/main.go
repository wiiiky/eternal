package main

import (
	"eternal/config"
	"eternal/event"
	"eternal/eworker/answer"
	"eternal/eworker/sms"
	"eternal/logging"
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
	"os"
)

const APPNAME = "eworker"

func main() {
	config.Init(APPNAME)
	initLogging()
	initDatabase()
	initEvent()
	initSMS()
	event.Register(event.KeyAnswerUpvote, answer.HandleAnswerUpvote)
	event.Register(event.KeyAnswerDownvote, answer.HandleAnswerDownvote)
	event.Register(event.KeySMSSend, sms.HandleSMSSend)

	ch := make(chan os.Signal)
	<-ch
}

func initEvent() {
	amqpURL := config.GetString("event.amqp.url")
	amqpExchange := config.GetString("event.amqp.exchange")
	amqpRouteKey := config.GetString("event.amqp.route_key")
	amqpQueue := config.GetString("event.amqp.queue")
	amqpConsumer := config.GetStringDefault("event.amqp.consumer", "eworker")
	event.InitSub(amqpURL, amqpExchange, amqpQueue, amqpRouteKey, amqpConsumer)
}

func initLogging() {
	format := config.GetStringDefault("log.format", "json")
	level := config.GetStringDefault("log.level", "info")
	output := config.GetStringDefault("log.output", "stdout")

	logging.Init(format, level, output)
}

func initDatabase() {
	pgURL := config.GetString("database.pg.url")
	mongoURL := config.GetString("database.mongo.url")
	mongoDBName := config.GetString("database.mongo.dbname")
	if pgURL == "" {
		log.Fatal("**CONFIG** database.pg.url not found")
	} else if mongoURL == "" {
		log.Fatal("**CONFIG** database.mongo.url not found")
	} else if mongoDBName == "" {
		log.Fatal("**CONFIG** database.mongo.dbname not found")
	}
	if err := db.Init(pgURL, mongoURL, mongoDBName); err != nil {
		log.Fatal("Connecting database failed:", err)
	}
}

func initSMS() {
	appid := config.GetString("sms.appid")
	appkey := config.GetString("sms.appkey")
	if appid == "" {
		log.Fatal("**CONFIG** sms.appid not found")
	} else if appkey == "" {
		log.Fatal("**CONFIG** sms.appkey not found")
	}
	keys := config.GetStringMapString("sms.keys")
	sms.Init(appid, appkey, keys)
}
