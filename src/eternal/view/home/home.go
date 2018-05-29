package home

import (
	"eternal/middleware"
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	"net/http"
)

func GetHotAnswers(ctx echo.Context) error {
	userID := ctx.Get(middleware.CTX_KEY_ACCOUNT_ID).(string)
	data := HotAnswerPageData{
		Before: "",
		Limit:  10,
	}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	if err := ctx.Validate(&data); err != nil {
		return err
	}
	answers, err := questionModel.FindHotAnswers(userID, data.Before, data.Limit)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, answers)
}
