package user

import (
	"eternal/controller/context"
	"eternal/errors"
	userModel "eternal/model/user"
	"github.com/labstack/echo"
	"net/http"
)

func GetUserProfile(c echo.Context) error {
	ctx := c.(*context.Context)
	up, err := userModel.GetUserProfile(ctx.Account.ID)
	if err != nil {
		return err
	} else if up == nil {
		return errors.ErrUserNotFound
	} else {
		return ctx.JSON(http.StatusOK, up)
	}
}

/* 更新用户的封面图 */
func UpdateUserCover(c echo.Context) error {
	ctx := c.(*context.Context)
	data := UpdateCoverRequest{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	if err := ctx.Validate(data); err != nil {
		return err
	}

	up, err := userModel.UpdateUserCover(ctx.Account.ID, data.Cover)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, up)
}
