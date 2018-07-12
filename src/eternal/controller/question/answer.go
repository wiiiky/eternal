package question

import (
	"eternal/controller/context"
	"eternal/event"
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/* 获取热门回答 */
func GetHotAnswers(c echo.Context) error {
	ctx := c.(*context.Context)
	data := HotAnswerPageRequest{
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

/* 获取问题下的回答 */
func GetQuestionAnswers(c echo.Context) error {
	ctx := c.(*context.Context)
	questionID := ctx.Param("qid")
	data := QuestionAnswerPageRequest{
		Page:  1,
		Limit: 10,
	}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	if err := ctx.Validate(&data); err != nil {
		return err
	}
	answers, err := questionModel.GetQuestionAnswers(ctx.Account.ID, questionID, data.Page, data.Limit)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, answers)
}

/* 点赞回答 */
func UpvoteAnswer(c echo.Context) error {
	ctx := c.(*context.Context)
	userID := ctx.Account.ID
	answerID := ctx.Param("id")
	upvoteCount, downvoteCount, err := questionModel.UpvoteAnswer(userID, answerID)
	if err != nil {
		return err
	}
	event.Publish(event.KeyAnswerUpvote, event.AnswerUpvote{
		AnswerID: answerID,
		UserID:   userID,
	})
	return ctx.JSON(http.StatusOK, &VoteAnswerResult{
		UpvoteCount:   upvoteCount,
		DownvoteCount: downvoteCount,
	})
}

func UndoUpvoteAnswer(c echo.Context) error {
	ctx := c.(*context.Context)
	userID := ctx.Account.ID
	answerID := ctx.Param("id")
	upvoteCount, downvoteCount, err := questionModel.UndoUpvoteAnswer(userID, answerID)
	if err != nil {
		return err
	}
	event.Publish(event.KeyAnswerDownvote, event.AnswerDownvote{
		AnswerID: answerID,
		UserID:   userID,
	})
	return ctx.JSON(http.StatusOK, &VoteAnswerResult{
		UpvoteCount:   upvoteCount,
		DownvoteCount: downvoteCount,
	})
}

/* 不赞同回答 */
func DownvoteAnswer(c echo.Context) error {
	ctx := c.(*context.Context)
	userID := ctx.Account.ID
	answerID := ctx.Param("id")
	upvoteCount, downvoteCount, err := questionModel.DownvoteAnswer(userID, answerID)
	if err != nil {
		return err
	}
	event.Publish(event.KeyAnswerDownvote, event.AnswerDownvote{
		AnswerID: answerID,
		UserID:   userID,
	})
	return ctx.JSON(http.StatusOK, &VoteAnswerResult{
		UpvoteCount:   upvoteCount,
		DownvoteCount: downvoteCount,
	})
}

func UndoDownvoteAnswer(c echo.Context) error {
	ctx := c.(*context.Context)
	userID := ctx.Account.ID
	answerID := ctx.Param("id")
	upvoteCount, downvoteCount, err := questionModel.UndoDownvoteAnswer(userID, answerID)
	if err != nil {
		return err
	}
	event.Publish(event.KeyAnswerUpvote, event.AnswerUpvote{
		AnswerID: answerID,
		UserID:   userID,
	})
	return ctx.JSON(http.StatusOK, &VoteAnswerResult{
		UpvoteCount:   upvoteCount,
		DownvoteCount: downvoteCount,
	})
}
