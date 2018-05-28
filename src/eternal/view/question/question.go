package question

import (
	"eternal/errors"
	"eternal/middleware"
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	"net/http"
)

func CreateQuestion(ctx echo.Context) error {
	userID := ctx.Get(middleware.CTX_KEY_ACCOUNT_ID).(string)
	data := CreateQuestionRequest{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	if err := ctx.Validate(data); err != nil {
		return err
	}

	question, err := questionModel.CreateQuestion(userID, data.Title, data.Topics, data.Content)
	if err != nil {
		return err
	} else if question == nil {
		return errors.ErrQuestionNotFound
	}
	return ctx.JSON(http.StatusOK, question)
}
