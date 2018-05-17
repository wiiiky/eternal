package answer

import (
	"eternal/model/account"
	"eternal/model/answer"
	"github.com/labstack/echo"
	"net/http"
)

func AddAnswerLike(ctx echo.Context) error {
	a := ctx.Get("account").(*account.Account)
	answerID := ctx.Param("id")
	err := answer.AddAnswerLike(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}

func AddAnswerDislike(ctx echo.Context) error {
	a := ctx.Get("account").(*account.Account)
	answerID := ctx.Param("id")
	err := answer.AddAnswerDislike(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
