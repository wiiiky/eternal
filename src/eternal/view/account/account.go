package account

import (
	accountModel "eternal/model/account"
	"eternal/model/db"
	"eternal/view/errors"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/* 获取当前支持的国家 */
func GetSupportedCountries(ctx echo.Context) error {
	countries, err := accountModel.GetSupportedCountries()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, countries)
}

func login(ctx echo.Context, a *accountModel.Account) error {
	tk, err := accountModel.UpsertToken(a.ID)
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
	data := LoginRequest{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	mobile := data.Mobile
	password := data.Password

	a, err := accountModel.GetAccountWithMobile(mobile)
	if err != nil {
		return err
	} else if a == nil { /* 用户不存在 */
		return errors.ErrUserNotFound
	}
	if !a.Auth(password) { /* 密码错误 */
		return errors.ErrUserPasswordInvalid
	}
	log.Debugf("User %s logged in", mobile)

	return login(ctx, a)
}

func Logout(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)

	if err := accountModel.DeleteToken(a.ID); err != nil {
		return err
	}

	sess, _ := session.Get("session", ctx)
	sess.Values["token"] = ""
	sess.Save(ctx.Request(), ctx.Response())
	return ctx.NoContent(http.StatusOK)
}

/* 注册 */
func Signup(ctx echo.Context) error {
	data := SignupRequest{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	countryCode := data.CountryCode
	mobile := data.Mobile
	password := data.Password
	code := data.Code

	if code != "123456" {
		return errors.ErrUseSMSCodeInvalid
	}
	if len(password) <= 4 || len(password) >= 12 {
		return errors.ErrUserPasswordLengthInvalid
	}

	country, err := accountModel.GetSupportedCountryWithCode(countryCode)
	if err != nil {
		return err
	} else if country == nil {
		return errors.ErrCountryCodeInvalid
	}

	a, err := accountModel.CreateAccount(countryCode, mobile, password, accountModel.PTYPE_MD5)
	if err == db.ErrKeyDuplicate {
		return errors.ErrMobileExisted
	} else if err != nil {
		return err
	}

	return login(ctx, a)
}

/* 获取帐号信息 */
func GetAccountInfo(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	return ctx.JSON(http.StatusOK, a)
}
