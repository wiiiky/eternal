package home

import (
	accountModel "eternal/model/account"
	questionModel "eternal/model/question"
	"eternal/view"
	"github.com/labstack/echo"
	"net/http"
)

func GetHotAnswers(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	var pd view.PageData
	if err := ctx.Bind(&pd); err != nil {
		return err
	}
	if pd.Page <= 0 {
		pd.Page = 1
	}
	if pd.Limit <= 0 {
		pd.Limit = 10
	}
	answers, err := questionModel.FindHotAnswers(a.ID, pd.Page, pd.Limit)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, answers)
}
