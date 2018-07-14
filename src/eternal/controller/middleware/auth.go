package middleware

import (
	accountCache "eternal/cache/account"
	"eternal/controller/context"
	"eternal/errors"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*context.Context)
		tokenID := ctx.GetCookieToken()
		if tokenID == "" {
			return c.NoContent(http.StatusUnauthorized)
		}
		token, err := accountCache.GetToken(tokenID)
		if err != nil {
			return err
		} else if token == nil {
			return c.NoContent(http.StatusUnauthorized)
		} else if token.ETime.Before(time.Now()) { /* TOKEN已过期 */
			return errors.ErrTokenExpired
		}
		account, _ := accountCache.GetAccount(token.UserID)
		if account == nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		ctx.Account = account
		return next(ctx)
	}
}
