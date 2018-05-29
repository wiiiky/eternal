package home

type HotAnswerPageData struct {
	Limit  int    `query:"limit" validate:"gte=1"`
	Before string `query:"before"` // 时间点
}
