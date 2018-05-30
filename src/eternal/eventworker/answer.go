package main

import (
	"encoding/json"
	"eternal/event"
	questionModel "eternal/model/question"
	log "github.com/sirupsen/logrus"
	"time"
)

/* 如果在一天内赞超过两次，则设置为热门回答 */
func handleAnswerUpvote(routeKey string, body []byte) bool {
	var data event.AnswerDownvote
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("json.Unmarshal failed:", err)
		return false
	}
	answerID := data.AnswerID
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1)
	count, err := questionModel.GetAnswerUpvoteCount(answerID, startTime, endTime)
	if err != nil {
		return false
	} else if count >= 2 {
		/* 设置为热门回答 */
		questionModel.UpsertHotAnswer(answerID)
	}
	return false
}

func handleAnswerDownvote(routeKey string, body []byte) bool {
	log.Infof("%s %s\n", routeKey, body)
	return false
}
