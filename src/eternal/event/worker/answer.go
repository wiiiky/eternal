package main

import (
	log "github.com/sirupsen/logrus"
)

func handleAnswerUpvote(routeKey string, body []byte) bool {
	log.Infof("%s %s\n", routeKey, body)
	return false
}

func handleAnswerDownvote(routeKey string, body []byte) bool {
	log.Infof("%s %s\n", routeKey, body)
	return false
}
