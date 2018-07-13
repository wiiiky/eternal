package sms

import (
	"eternal/errors"
	"eternal/model/db"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"time"
)

/* 获取指定手机号码的验证码 */
func FindSMSCode(phoneNumber, codeType string, codeStatus int, d time.Duration) (*SMSCode, error) {
	mc := db.MC(CollectionSMSCode)
	smsCode := SMSCode{}
	etime := time.Now().Add(-d)

	query := bson.M{
		"phone_number": phoneNumber,
		"type":         codeType,
		"status":       codeStatus,
		"ctime": bson.M{
			"$gt": etime,
		},
	}
	if err := mc.Find(query).Sort("-ctime").One(&smsCode); err != nil {
		if err != db.ErrNotFound {
			log.Error("MGO Error:", err)
			return nil, errors.ErrDB
		}
		return nil, nil
	}
	return &smsCode, nil
}

/* 获取指定IP在过去一段时间内的短信发送数量 */
func CountSMSCodeByClientIP(clientIP, codeType string, d time.Duration) (int, error) {
	mc := db.MC(CollectionSMSCode)
	etime := time.Now().Add(-d)
	query := bson.M{
		"client_ip": clientIP,
		"type":      codeType,
		"ctime": bson.M{
			"$gt": etime,
		},
	}

	if count, err := mc.Find(query).Count(); err != nil {
		log.Error("MGO Error:", err)
		return 0, errors.ErrDB
	} else {
		return count, nil
	}
}

/* 更新短信验证码状态 */
func UpdateSMSCodeStatus(ID bson.ObjectId, codeStatus int) error {
	mc := db.MC(CollectionSMSCode)
	query := bson.M{
		"_id": ID,
	}
	change := bson.M{
		"$set": bson.M{
			"status": codeStatus,
			"utime":  time.Now(),
		},
	}
	if err := mc.Update(query, change); err != nil {
		log.Error("MGO Error:", err)
		return errors.ErrDB
	}
	return nil
}

/* 插入短信验证码 */
func InsertSMSCode(phoneNumber, codeType, code, clientIP string, codeStatus int) (*SMSCode, error) {
	mc := db.MC(CollectionSMSCode)
	smsCode := SMSCode{
		PhoneNumber: phoneNumber,
		Code:        code,
		Type:        codeType,
		Status:      codeStatus,
		ClientIP:    clientIP,
		UTime:       time.Now(),
		CTime:       time.Now(),
	}
	if err := mc.Insert(&smsCode); err != nil {
		log.Error("MGO Error:", err)
		return nil, errors.ErrDB
	}
	return &smsCode, nil
}
