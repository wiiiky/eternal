package user

import (
	"eternal/model/account"
	"eternal/model/user"
	"eternal/view/errors"
	"github.com/labstack/echo"
	"net/http"
)

func GetUserProfile(ctx echo.Context) error {
	a := ctx.Get("account").(*account.Account)
	up, err := user.GetUserProfile(a.ID)
	if err != nil {
		return err
	} else if up == nil {
		return errors.ErrUserNotFound
	} else {
		return ctx.JSON(http.StatusOK, up)
	}
}
