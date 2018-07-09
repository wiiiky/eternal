package question

import (
	"eternal/controller/context"
	"eternal/errors"
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

/* 创建问题 */
func CreateQuestion(c echo.Context) error {
	ctx := c.(*context.Context)
	data := CreateQuestionRequest{}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	if err := ctx.Validate(data); err != nil {
		return err
	}

	question, err := questionModel.CreateQuestion(ctx.Account.ID, data.Title, data.Topics, data.Content)
	if err != nil {
		return err
	} else if question == nil {
		return errors.ErrQuestionNotFound
	}
	return ctx.JSON(http.StatusOK, question)
}

/* 关注问题 */
func FollowQuestion(c echo.Context) error {
	ctx := c.(*context.Context)
	questionID := ctx.Param("id")
	
	if followCount, err := questionModel.FollowQuestion(ctx.Account.ID, questionID); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, FollowQuestionResult{
			FollowCount: followCount,
		})
	}
}

/* 取消关注 */
func UnfollowQuestion(c echo.Context) error {
	ctx := c.(*context.Context)
	questionID := ctx.Param("id")
	if followCount, err := questionModel.UnfollowQuestion(ctx.Account.ID, questionID); err != nil {
		return err
	} else {
		return ctx.JSON(http.StatusOK, FollowQuestionResult{
			FollowCount: followCount,
		})
	}
}

/* 搜索问题 */
func FindQuestions(c echo.Context) error {
	ctx := c.(*context.Context)
	userID := ctx.Account.ID
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
