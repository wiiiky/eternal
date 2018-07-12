package question

import (
	questionModel "eternal/model/question"
)

type HotAnswerPageRequest struct {
	Limit  int    `query:"limit" validate:"gte=1"`
	Before string `query:"before"` // 时间点
}

type HotAnswer struct {
	*questionModel.HotAnswer
	UserAnswerRelationship *questionModel.UserAnswerRelationship `json:"user_answer_relationship"`
}

type VoteAnswerResult struct {
	UpvoteCount   uint64 `json:"upvote_count"`
	DownvoteCount uint64 `json:"downvote_count"`
}

type SearchTopicRequest struct {
	Query string `query:"query" validate:"required"`
	Limit int    `query:"limit" validate:"gte=1"`
}

type SearchQuestionRequest struct {
	Query string `query:"query" validate:"required"`
	Limit int    `query:"limit" validate:"gte=1"`
	Page  int    `query:"page" validate:"gte=1"`
}

type CreateQuestionRequest struct {
	Title   string   `json:"title" form:"title" validate:"required"`
	Topics  []string `json:"topics" form:"topics" validate:"gt=0,required"`
	Content string   `json:"content" form:"content"`
}

type QuestionAnswerPageRequest struct {
	Limit int `query:"limit" validate:"gte=1"`
	Page  int `query:"page" validate:"gte=1"`
}

type QuestionWithTopAnswer struct {
	Question               *questionModel.Question               `json:"question"`
	TopAnswer              *questionModel.Answer                 `json:"answer"`
	UserAnswerRelationship *questionModel.UserAnswerRelationship `json:"user_answer_relationship"`
}

type FollowQuestionResult struct {
	FollowCount uint64 `json:"follow_count"`
	Followed    bool   `json:"followed"`
}

type QuestionResult struct {
	*questionModel.Question
	UserQuestionRelationship *questionModel.UserQuestionRelationship `json:"user_question_relationship"`
}
