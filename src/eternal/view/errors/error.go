package errors

import (
	"fmt"
	"net/http"
)

/* 这里的错误全部都是对客户端的 */
var (
	ErrUserNotFound              = NewError(http.StatusNotFound, 100, "用户不存在")
	ErrUserPasswordInvalid       = NewError(http.StatusBadRequest, 101, "用户密码错误")
	ErrUserPasswordLengthInvalid = NewError(http.StatusBadRequest, 102, "账号密码长度过短或过长")
	ErrUseSMSCodeInvalid         = NewError(http.StatusBadRequest, 103, "短信验证码错误")
	ErrMobileExisted             = NewError(http.StatusBadRequest, 104, "用户已存在")
	ErrCountryCodeInvalid        = NewError(http.StatusBadRequest, 105, "国家不支持")

	ErrFileNotFound = NewError(http.StatusNotFound, 1001, "文件不存在")
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
