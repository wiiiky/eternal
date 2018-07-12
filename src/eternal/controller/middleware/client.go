package middleware

import (
	clientCache "eternal/cache/client"
	"eternal/controller/context"
	"eternal/errors"
	"github.com/labstack/echo"
)

/* 验证客户端信息 */
func SourceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		clientID := c.Request().Header.Get("Source")
		if clientID == "" {
			clientID = c.QueryParam("_s")
			if clientID == "" {
				return errors.ErrClientInvalid
			}
		}
		client, err := clientCache.GetClient(clientID)
		if err != nil {
			return err
		} else if client == nil {
			return errors.ErrClientInvalid
		}
		return next(&context.Context{
			Context: c,
			Client:  client,
		})
	}
}
