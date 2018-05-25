package question

import (
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	"net/http"
)

func FindTopics(ctx echo.Context) error {
	data := SearchTopicRequest{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	if err := ctx.Validate(data); err != nil {
		return err
	}
	topics, err := questionModel.FindTopics(data.Query, data.Limit)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, topics)
}
