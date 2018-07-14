package errors

import (
	"fmt"
	"net/http"
)

/* 这里的错误全部都是对客户端的 */
var (
	ErrDB            = NewError(http.StatusInternalServerError, -1, "数据库错误")
	ErrClientInvalid = NewError(http.StatusBadRequest, -2, "客户端未定义")

	/* 用户相关 */
	ErrUserNotFound              = NewError(http.StatusNotFound, 100, "用户不存在")
	ErrUserPasswordInvalid       = NewError(http.StatusBadRequest, 101, "用户密码错误")
	ErrUserPasswordLengthInvalid = NewError(http.StatusBadRequest, 102, "账号密码长度过短或过长")
	ErrUseSMSCodeInvalid         = NewError(http.StatusBadRequest, 103, "短信验证码错误")
	ErrPhoneNumberExisted        = NewError(http.StatusBadRequest, 104, "手机号已存在")
	ErrCountryCodeInvalid        = NewError(http.StatusBadRequest, 105, "国家不支持")
	ErrTokenExpired              = NewError(http.StatusForbidden, 106, "token过期")

	/* 11 MISC */
	ErrSMSTooOften = NewError(http.StatusBadRequest, 1101, "短信发送太快")
	ErrSMSTooMany  = NewError(http.StatusBadRequest, 1102, "短时间内发送过多")

	/* 各种不存在 */
	ErrFileNotFound     = NewError(http.StatusNotFound, 1001, "文件不存在")
	ErrQuestionNotFound = NewError(http.StatusNotFound, 1002, "问题不存在")
	ErrTopicNotFound    = NewError(http.StatusNotFound, 1003, "话题不存在")
	ErrAnswerNotFound   = NewError(http.StatusNotFound, 1004, "回答不存在")
)

func NewError(status, code int, msg interface{}) *Error {
	return &Error{
		HttpStatus: status,
		Code:       code,
		Message:    fmt.Sprintf("%v", msg),
	}
}

func CopyError(e *Error) *Error {
	err := new(Error)
	*err = *e
	return err
}

func CopyErrorWithMsg(e *Error, msg string) *Error {
	err := CopyError(e)
	err.Message = msg
	return err
}

type Error struct {
	HttpStatus int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
