package user

import (
	"eternal/errors"
	accountModel "eternal/model/account"
	userModel "eternal/model/user"
	"github.com/labstack/echo"
	"net/http"
)

func GetUserProfile(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	up, err := userModel.GetUserProfile(a.ID)
	if err != nil {
		return err
	} else if up == nil {
		return errors.ErrUserNotFound
	} else {
		return ctx.JSON(http.StatusOK, up)
	}
}

/* 更新用户的封面图 */
func UpdateUserCover(ctx echo.Context) error {
	data := UpdateCoverRequest{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	if err := ctx.Validate(data); err != nil {
		return err
	}
	a := ctx.Get("account").(*accountModel.Account)

	up, err := userModel.UpdateUserCover(a.ID, data.Cover)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, up)
}
