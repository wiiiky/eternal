package question

import (
	accountModel "eternal/model/account"
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	"net/http"
)

func UpvoteAnswer(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	err := questionModel.UpvoteAnswer(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}

func DownvoteAnswer(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	err := questionModel.DownvoteAnswer(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
