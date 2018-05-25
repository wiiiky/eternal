package question

type VoteAnswerResult struct {
	UpvoteCount   uint64 `json:"upvote_count"`
	DownvoteCount uint64 `json:"downvote_count"`
}

type SearchTopicRequest struct {
	Query string `query:"q" validate:"required"`
	Limit int    `query:"limit validate:"default=10"`
}