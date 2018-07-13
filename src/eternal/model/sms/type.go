package sms

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

/* mongodb的Collection */
const (
	CollectionSMSCode = "sms_code"
)

/* 短信验证码状态 */
const (
	CodeStatusUnused = 0
	CodeStatusUsed   = 1
)

/* 短信验证码类型 */
const (
	CodeTypeSignup = "signup"
)

type SMSCode struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	PhoneNumber string        `json:"phone_number" bson:"phone_number"`
	Code        string        `json:"-" bson:"code"`
	Status      int           `json:"status" bson:"status"`
	Type        string        `json:"type" bson:"type"`
	ClientIP    string        `json:"client_ip" bson:"client_ip"`
	UTime       time.Time     `json:"utime" bson:"utime"`
	CTime       time.Time     `json:"ctime" bson:"ctime"`
}
