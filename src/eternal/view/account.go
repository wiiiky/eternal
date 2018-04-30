package view

import (
	"eternal/model/account"
	"eternal/model/db"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/* 获取当前支持的国家 */
func GetSupportedCountries(ctx echo.Context) error {
	countries, err := account.GetSupportedCountries()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, countries)
}

type iLogin struct {
	CountryCode string `json:"country_code" form:"country_code" query:"country_code"`
	Mobile      string `json:"mobile" form:"mobile" query:"mobile"`
	Password    string `json:"password" form:"password" query:"password"`
}

func login(ctx echo.Context, a *account.Account) error {
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

/* 登录 */
func Login(ctx echo.Context) error {
	data := iLogin{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	countryCode := data.CountryCode
	mobile := data.Mobile
	password := data.Password

	a, err := account.GetAccountWithMobile(countryCode, mobile)
	if err != nil {
		return err
	} else if a == nil { /* 用户不存在 */
		return ErrUserNotFound
	}
	if !a.Auth(password) { /* 密码错误 */
		return ErrUserPasswordInvalid
	}
	log.Debugf("User %s:%s logged in", countryCode, mobile)

	return login(ctx, a)
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
	countryCode := data.CountryCode
	mobile := data.Mobile
	password := data.Password
	code := data.Code

	if code != "123456" {
		return ErrUseSMSCodeInvalid
	}
	if len(password) <= 4 || len(password) >= 12 {
		return ErrUserPasswordLengthInvalid
	}

	country, err := account.GetSupportedCountryWithCode(countryCode)
	if err != nil {
		return err
	} else if country == nil {

	}

	a, err := account.CreateAccount(countryCode, mobile, password, account.PTYPE_MD5)
	if err == db.ErrKeyDuplicate {
		return ErrMobileExisted
	} else if err != nil {
		return err
	}

	return login(ctx, a)
}

/* 获取帐号信息 */
func GetAccountInfo(ctx echo.Context) error {
	a := ctx.Get("account").(*account.Account)
	return ctx.JSON(http.StatusOK, a)
}
