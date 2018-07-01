package home

import (
	"eternal/controller/context"
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetHotAnswers(c echo.Context) error {
	ctx := c.(*context.Context)
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
	userID := ctx.Account.ID
	hotAnswers, err := questionModel.FindHotAnswers(userID, data.Before, data.Limit)
	if err != nil {
		return err
	}
	results := make([]*HotAnswer, 0)
	for _, hotAnswer := range hotAnswers {
		relationship, err := questionModel.GetUserAnswerRelationship(userID, hotAnswer.Answer.ID)
		if err != nil {
			log.Error("GetUserAnswerRelationship failed:", err)
			return err
		}
		results = append(results, &HotAnswer{HotAnswer: hotAnswer, UserAnswerRelationship: relationship})
	}
	return ctx.JSON(http.StatusOK, results)
}
