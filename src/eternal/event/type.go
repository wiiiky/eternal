package event

const (
	KeyAnswerUpvote   = "answer.upvote"
	KeyAnswerDownvote = "answer.downvote"
	KeySMSSend        = "sms.send"
)

type AnswerUpvoteData struct {
	AnswerID string `json:"answer_id"`
	UserID   string `json:"user_id"`
}

type AnswerDownvoteData struct {
	AnswerID string `json:"answer_id"`
	UserID   string `json:"user_id"`
}

const (
	SMSKeySignup = "signup"
)

type SMSSendData struct {
	PhoneNumber string            `json:"phone_number"`
	Vars        map[string]string `json:"vars"`
	Key         string            `json:"key"`
}
