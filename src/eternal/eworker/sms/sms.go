package sms

import (
	"encoding/json"
	"eternal/event"
	"eternal/eworker/sms/submail"
	log "github.com/sirupsen/logrus"
)

var smsKeys map[string]string

func Init(appid, appkey string, keys map[string]string) error {
	smsKeys = keys
	if smsKeys == nil {
		smsKeys = make(map[string]string)
	}
	return submail.Init(appid, appkey)
}

/* 发送短信 */
func HandleSMSSend(routeKey string, body []byte) bool {
	var data event.SMSSendData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Error("json.Unmarshal failed:", err)
		return false
	}
	template := smsKeys[data.Key]
	if template == "" {
		log.Error("SMS key not found:", data.Key)
		return false
	}
	result, err := submail.XSend(data.PhoneNumber, template, data.Vars)
	if err != nil {
		log.Error("XSend failed:", err)
		return false
	} else if result.Status != submail.StatusSuccess {
		log.Error("XSend result error:", result.Status, result.Msg)
		return false
	} else {
		log.Info("XSend successfully:", data.PhoneNumber, data.Key)
	}
	return false
}
