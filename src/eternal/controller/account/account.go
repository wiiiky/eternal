package account

import (
	accountCache "eternal/cache/account"
	"eternal/controller/context"
	"eternal/errors"
	accountModel "eternal/model/account"
	smsModel "eternal/model/sms"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

/* 获取当前支持的国家 */
func GetSupportedCountries(ctx echo.Context) error {
	countries, err := accountModel.GetSupportedCountries()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, countries)
}

func login(ctx *context.Context, a *accountModel.Account) error {
	client := ctx.Client
	maxAge := client.TokenMaxAge // 单位秒

	/* 更新Token */
	tk, err := accountModel.UpsertToken(a.ID, ctx.Client.ID, time.Second*time.Duration(maxAge))
	if err != nil {
		return err
	}

	sess, _ := session.Get("session", ctx)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: int(maxAge),
	}
	sess.Values["token"] = tk.ID
	sess.Save(ctx.Request(), ctx.Response())
	return ctx.JSON(http.StatusOK, tk)
}

/* 登录 */
func Login(c echo.Context) error {
	ctx := c.(*context.Context)
	data := LoginRequest{}
	if err := ctx.Bind(&data); err != nil {
		return err
	} else if err := ctx.Validate(&data); err != nil {
		return err
	}
	phoneNumber := data.PhoneNumber
	password := data.Password

	a, err := accountModel.GetAccountByPhoneNumber(phoneNumber)
	if err != nil {
		return err
	} else if a == nil { /* 用户不存在 */
		return errors.ErrUserNotFound
	}
	if !a.Auth(password) { /* 密码错误 */
		return errors.ErrUserPasswordInvalid
	}
	log.Debugf("User %s logged in", phoneNumber)

	return login(ctx, a)
}

func Logout(c echo.Context) error {
	ctx := c.(*context.Context)

	if err := accountCache.DeleteToken(ctx.Token.ID); err != nil {
		return err
	}

	sess, _ := session.Get("session", ctx)
	sess.Values["token"] = ""
	sess.Save(ctx.Request(), ctx.Response())
	return ctx.NoContent(http.StatusOK)
}

/* 注册 */
func Signup(c echo.Context) error {
	ctx := c.(*context.Context)
	data := SignupRequest{}
	if err := ctx.Bind(&data); err != nil {
		return err
	} else if err := ctx.Validate(&data); err != nil {
		return err
	}
	countryCode := data.CountryCode
	phoneNumber := data.PhoneNumber
	password := data.Password
	code := data.Code

	if len(password) <= 4 || len(password) >= 12 {
		return errors.ErrUserPasswordLengthInvalid
	}

	country, err := accountModel.GetSupportedCountryByCode(countryCode)
	if err != nil {
		return err
	} else if country == nil {
		return errors.ErrCountryCodeInvalid
	}

	/* 验证短信验证码，CheckSMSCode方法会设置验证码为已使用 */
	if ok, err := smsModel.CheckSMSCode(phoneNumber, smsModel.CodeTypeSignup, code, time.Minute*20); err != nil {
		return err
	} else if !ok {
		return errors.ErrUseSMSCodeInvalid
	}

	a, err := accountModel.CreateAccount(countryCode, phoneNumber, password, accountModel.PTYPE_MD5)
	if err != nil {
		return err
	}

	return login(ctx, a)
}

/* 获取帐号信息 */
func GetAccountInfo(c echo.Context) error {
	ctx := c.(*context.Context)
	return ctx.JSON(http.StatusOK, ctx.Account)
}
