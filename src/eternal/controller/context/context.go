package context

import (
	"github.com/labstack/echo-contrib/session"
	"strings"
)

func (ctx *Context) getSessionToken() string {
	sess, _ := session.Get("session", ctx)
	if sess == nil {
		return ""
	}
	v, ok := sess.Values["token"]
	if !ok || v == nil {
		return ""
	}
	return v.(string)
}

func (ctx *Context) getBearerToken() string {
	auth := ctx.Request().Header.Get("Authorization")
	segments := strings.Fields(auth)
	if len(segments) != 2 || strings.ToLower(segments[0]) != "bearer" {
		return ""
	}
	return segments[1]
}

func (ctx *Context) GetCookieToken() string {
	token := ctx.getSessionToken()
	if token == "" {
		token = ctx.getBearerToken()
	}
	return token
}
