package view

import (
	"fmt"
	"net/http"
)

/* 这里的错误全部都是对客户端的 */
var (
	E_USER_NOT_FOUND               = NewError(http.StatusNotFound, 100, "用户不存在")
	E_USER_PASSWORD_INVALID        = NewError(http.StatusBadRequest, 101, "用户密码错误")
	E_USER_PASSWORD_LENGTH_INVALID = NewError(http.StatusBadRequest, 102, "账号密码长度过短或过长")
	E_USER_SMS_CODE_INVALID        = NewError(http.StatusBadRequest, 103, "短信验证码错误")
)

func NewError(status, code int, msg interface{}) error {
	return &Error{
		HttpStatus: status,
		Code:       code,
		Message:    fmt.Sprintf("%v", msg),
	}
}

type Error struct {
	HttpStatus int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
