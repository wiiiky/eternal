package submail

import (
	"encoding/json"
)

var messageconfig = make(map[string]string)

func Init(appid, appkey string) error {
	messageconfig["appid"] = appid
	messageconfig["appkey"] = appkey
	messageconfig["signtype"] = "md5"

	return nil
}

func XSend(phoneNumber, template string, vars map[string]string) (*XSendResult, error) {
	messagexsend := CreateMessageXSend()
	MessageXSendAddTo(messagexsend, phoneNumber)
	MessageXSendSetProject(messagexsend, template)
	MessageXSendAddVars(messagexsend, vars)
	data := MessageXSendRun(MessageXSendBuildRequest(messagexsend), messageconfig)

	result := XSendResult{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
