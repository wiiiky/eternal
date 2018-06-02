package question

import (
	"eternal/errors"
	"eternal/middleware"
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	"net/http"
)

/* 获取问题的详细信息 */
func GetQuestion(ctx echo.Context) error {
	questionID := ctx.Param("id")
	question, err := questionModel.GetQuestion(questionID)
	if err != nil {
		return err
	} else if question == nil {
		return errors.ErrQuestionNotFound
	}
	return ctx.JSON(http.StatusOK, question)
}

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

/* 搜索问题 */
func FindQuestions(ctx echo.Context) error {
	userID := ctx.Get(middleware.CTX_KEY_ACCOUNT_ID).(string)
	data := SearchQuestionRequest{
		Page:  1,
		Limit: 10,
	}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	if err := ctx.Validate(&data); err != nil {
		return err
	}
	questions, err := questionModel.FindQuestions(data.Query, data.Page, data.Limit)
	if err != nil {
		return err
	}
	results := make([]*QuestionWithTopAnswer, 0)
	for _, question := range questions {
		topAnswer, err := questionModel.GetQuestionTopAnswer(question.ID)
		if err != nil {
			return err
		}
		var relationship *questionModel.UserAnswerRelationship
		if topAnswer != nil {
			relationship, err = questionModel.GetUserAnswerRelationship(userID, topAnswer.ID)
		}
		results = append(results, &QuestionWithTopAnswer{
			Question:               question,
			TopAnswer:              topAnswer,
			UserAnswerRelationship: relationship,
		})
	}
	return ctx.JSON(http.StatusOK, results)
}
