package sms

import (
	"encoding/json"
	"eternal/event"
	log "github.com/sirupsen/logrus"
)

var smsKeys map[string]string

func Init(keys map[string]string) error {
	smsKeys = keys
	if smsKeys == nil {
		smsKeys = make(map[string]string)
	}
	return nil
}

/* 如果在一天内赞超过两次，则设置为热门回答 */
func HandleSMSSend(routeKey string, body []byte) bool {
	var data event.SMSSendData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("json.Unmarshal failed:", err)
		return false
	}
	log.Info(data)
	return false
}
