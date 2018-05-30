package event

const (
	KeyAnswerUpvote   = "answer.upvote"
	KeyAnswerDownvote = "answer.downvote"
)

type AnswerUpvote struct {
	AnswerID string `json:"answer_id"`
	UserID   string `json:"user_id"`
}

type AnswerDownvote struct {
	AnswerID string `json:"answer_id"`
	UserID   string `json:"user_id"`
}
