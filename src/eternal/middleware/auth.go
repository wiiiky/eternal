package middleware

import (
	"eternal/model/account"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	log "github.com/sirupsen/logrus"
	"net/http"
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
		log.Debugf("Path = %s Session = %v", c.Path(), sess)
		if sess == nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		v, ok := sess.Values["token"]
		if !ok {
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
		c.Set("account", account)
		return next(c)
	}
}
