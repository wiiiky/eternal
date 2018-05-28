package middleware

import (
	"eternal/model/account"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	CTX_KEY_ACCOUNT    = "account"
	CTX_KEY_ACCOUNT_ID = "account.id"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	var freePathes = map[string]bool{
		"/login":  true,
		"/signup": true,
	}
	return func(c echo.Context) error {
		path := c.Path()
		if v, ok := freePathes[path]; v && ok {
			return next(c)
		}

		sess, _ := session.Get("session", c)
		if sess == nil {
			log.Debugf("session not found")
			return c.NoContent(http.StatusUnauthorized)
		}
		v, ok := sess.Values["token"]
		if !ok {
			log.Debugf("token not found %v", sess)
			return c.NoContent(http.StatusUnauthorized)
		}
		tokenID, ok := v.(string)
		if !ok {
			return c.NoContent(http.StatusUnauthorized)
		}
		account, _ := account.GetAccountWithTokenID(tokenID)
		if account == nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		c.Set(CTX_KEY_ACCOUNT, account)
		c.Set(CTX_KEY_ACCOUNT_ID, account.ID)
		return next(c)
	}
}
