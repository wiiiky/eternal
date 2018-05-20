package question

import (
	accountModel "eternal/model/account"
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	"net/http"
)

func AddAnswerLike(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	err := questionModel.AddAnswerLike(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}

func AddAnswerDislike(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	err := questionModel.AddAnswerDislike(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
