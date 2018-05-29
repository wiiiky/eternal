package view

type PageData struct {
	Page  int `query:"page" validate:"gte=1"`
	Limit int `query:"limit" validate:"gte=1"`
	After int `query:"after"` // 时间点
}
