package question

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

type QuestionAnswerPageData struct {
	Limit int `query:"limit" validate:"gte=1"`
	Page  int `query:"page" validate:"gte=1"`
}
