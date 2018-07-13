package misc

import (
	"eternal/controller/context"
	"eternal/errors"
	"eternal/event"
	smsModel "eternal/model/sms"
	"eternal/util"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

/* 发送注册验证码短信 */
func SendSignupCode(c echo.Context) error {
	ctx := c.(*context.Context)
	var data SendSignupCodeRequest
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	phoneNumber := data.PhoneNumber
	clientIP := ctx.RealIP()

	/* 如果在一个小时内发送短信超过5条，则返回错误 */
	if count, err := smsModel.CountSMSCodeByClientIP(clientIP, smsModel.CodeTypeSignup, time.Hour); err != nil {
		return err
	} else if count >= 5 {
		return errors.ErrSMSTooMany
	}

	smsCode, err := smsModel.FindSMSCode(phoneNumber, smsModel.CodeTypeSignup, smsModel.CodeStatusUnused, time.Minute)
	if err != nil {
		return err
	} else if smsCode != nil {
		return ctx.JSON(http.StatusOK, &SendSignupCodeResult{
			Sent: false,
			Wait: int(time.Minute/time.Second - time.Now().Sub(smsCode.CTime)/time.Second),
		})
	}
	code := util.RandDigit(6)
	smsCode, err = smsModel.InsertSMSCode(phoneNumber, smsModel.CodeTypeSignup, code, clientIP, smsModel.CodeStatusUnused)
	if err != nil {
		return err
	}
	/* 推送发送短信的验证码 */
	event.Publish(event.KeySMSSend, event.SMSSendData{
		PhoneNumber: phoneNumber,
		Key:         event.SMSKeySignup,
		Vars: map[string]string{
			"code": code,
		},
	})
	return ctx.JSON(http.StatusOK, &SendSignupCodeResult{
		Sent: true,
		Wait: int(time.Minute / time.Second),
	})
}
