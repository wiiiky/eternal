package question

import (
	accountModel "eternal/model/account"
	questionModel "eternal/model/question"
	"github.com/labstack/echo"
	"net/http"
)

func UpvoteAnswer(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	upvoteCount, downvoteCount, err := questionModel.UpvoteAnswer(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, &VoteAnswerResult{
		UpvoteCount:   upvoteCount,
		DownvoteCount: downvoteCount,
	})
}

func UndoUpvoteAnswer(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	upvoteCount, downvoteCount, err := questionModel.UndoUpvoteAnswer(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, &VoteAnswerResult{
		UpvoteCount:   upvoteCount,
		DownvoteCount: downvoteCount,
	})
}

func DownvoteAnswer(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	upvoteCount, downvoteCount, err := questionModel.DownvoteAnswer(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, &VoteAnswerResult{
		UpvoteCount:   upvoteCount,
		DownvoteCount: downvoteCount,
	})
}

func UndoDownvoteAnswer(ctx echo.Context) error {
	a := ctx.Get("account").(*accountModel.Account)
	answerID := ctx.Param("id")
	upvoteCount, downvoteCount, err := questionModel.UndoDownvoteAnswer(a.ID, answerID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, &VoteAnswerResult{
		UpvoteCount:   upvoteCount,
		DownvoteCount: downvoteCount,
	})
}
