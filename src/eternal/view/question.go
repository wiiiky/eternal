package view

import (
	"eternal/model/account"
	"eternal/model/question"
	"github.com/labstack/echo"
	"net/http"
)

func FindQuestions(ctx echo.Context) error {
	a := ctx.Get("account").(*account.Account)
	questions, err := question.FindQuestions(a.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, questions)
}

/* 获取热门问题 */
func FindHotQuestions(ctx echo.Context) error {
	return nil
}
