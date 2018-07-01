package middleware

import (
	"eternal/controller/context"
	accountModel "eternal/model/account"
	"github.com/labstack/echo"
	"net/http"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*context.Context)
		tokenID := ctx.GetCookieToken()
		if tokenID == "" {
			return c.NoContent(http.StatusUnauthorized)
		}
		account, _ := accountModel.GetAccountWithTokenID(tokenID)
		if account == nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		ctx.Account = account
		return next(ctx)
	}
}
