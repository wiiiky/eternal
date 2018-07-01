package home

import (
	questionModel "eternal/model/question"
)

type HotAnswerPageData struct {
	Limit  int    `query:"limit" validate:"gte=1"`
	Before string `query:"before"` // 时间点
}

type HotAnswer struct {
	*questionModel.HotAnswer
	UserAnswerRelationship *questionModel.UserAnswerRelationship `sql:"-" json:"user_answer_relationship"`
}
