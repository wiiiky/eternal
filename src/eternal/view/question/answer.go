package question

import (
	"eternal/event"
	"eternal/middleware"
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	"net/http"
)

/* 获取问题下的回答 */
func GetQuestionAnswers(ctx echo.Context) error {
	userID := ctx.Get(middleware.CTX_KEY_ACCOUNT_ID).(string)
	questionID := ctx.Param("qid")
	data := QuestionAnswerPageData{
		Page:  1,
		Limit: 10,
	}
	if err := ctx.Bind(&data); err != nil {
		return err
	}
	if err := ctx.Validate(&data); err != nil {
		return err
	}
	answers, err := questionModel.GetQuestionAnswers(userID, questionID, data.Page, data.Limit)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, answers)
}

func UpvoteAnswer(ctx echo.Context) error {
	userID := ctx.Get(middleware.CTX_KEY_ACCOUNT_ID).(string)
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

func UndoUpvoteAnswer(ctx echo.Context) error {
	userID := ctx.Get(middleware.CTX_KEY_ACCOUNT_ID).(string)
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

func DownvoteAnswer(ctx echo.Context) error {
	userID := ctx.Get(middleware.CTX_KEY_ACCOUNT_ID).(string)
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

func UndoDownvoteAnswer(ctx echo.Context) error {
	userID := ctx.Get(middleware.CTX_KEY_ACCOUNT_ID).(string)
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
