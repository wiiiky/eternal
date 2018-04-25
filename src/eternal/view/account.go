package view

import (
	"eternal/model/account"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type iLogin struct {
	CountryCode string `json:"country_code" form:"country_code" query:"country_code"`
	Mobile      string `json:"mobile" form:"mobile" query:"mobile"`
	Password    string `json:"password" form:"password" query:"password"`
}

/* 登录 */
func Login(ctx echo.Context) error {
	data := iLogin{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	countryCode := data.CountryCode
	mobile := data.Mobile
	password := data.Password

	a, err := account.GetWithMobile(countryCode, mobile)
	if err != nil {
		return err
	} else if a == nil { /* 用户不存在 */
		return E_USER_NOT_FOUND
	}
	if !a.Auth(password) { /* 密码错误 */
		return E_USER_PASSWORD_INVALID
	}
	log.Debugf("User %s:%s logged in", countryCode, mobile)

	tk, err := account.UpsertToken(a.ID)
	if err != nil {
		return err
	}

	sess, _ := session.Get("session", ctx)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	}
	sess.Values["token"] = tk.ID
	sess.Save(ctx.Request(), ctx.Response())
	return ctx.JSON(http.StatusOK, tk)
}

type iSignup struct {
	iLogin
	Code string `json:"code" form:"code" query:"code"`
}

/* 注册 */
func Signup(ctx echo.Context) error {
	data := iSignup{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	// countryCode := data.CountryCode
	// mobile := data.Mobile
	password := data.Password
	code := data.Code

	if code != "123456" {
		return E_USER_SMS_CODE_INVALID
	}
	if len(password) <= 4 || len(password) >= 12 {
		return E_USER_PASSWORD_LENGTH_INVALID
	}
	return nil
}

func GetAccountInfo(ctx echo.Context) error {
	a := ctx.Get("account").(*account.Account)
	return ctx.JSON(http.StatusOK, a)
}
