package answer

import (
	accountModel "eternal/model/account"
	answerModel "eternal/model/answer"
	"github.com/labstack/echo"
	"net/http"
)

func AddAnswerLike(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	err := answerModel.AddAnswerLike(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}

func AddAnswerDislike(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	err := answerModel.AddAnswerDislike(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
