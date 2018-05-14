package answer

import (
	"time"
)

type Answer struct {
	TableName struct{}  `sql:"answer" json:"answer"`
	ID        string    `sql:"id" json:"id"`
	Content   string    `sql:"content" json:"content"`
	UserID    string    `sql:"user_id" json:"user_id"`
	UTime     time.Time `sql:"utime,null" json:"utime"`
	CTime     time.Time `sql:"ctime,null" json:"ctime"`
}

type AnswerLike struct {
	TableName struct{}  `sql:"answer_like" json:"-"`
	UserID    string    `sql:"user_id,pk" json:"answer_id"`
	AnswerID  string    `sql:"answer_id,pk" json:"answer_id"`
	CTime     time.Time `sql:"ctime,null" json:"ctime"`
}

type AnswerDislike struct {
	TableName struct{}  `sql:"answer_dislike" json:"-"`
	UserID    string    `sql:"user_id,pk" json:"answer_id"`
	AnswerID  string    `sql:"answer_id,pk" json:"answer_id"`
	CTime     time.Time `sql:"ctime,null" json:"ctime"`
}

type AnswerStat struct {
	TableName    struct{} `sql:"answer_stat" json:"-"`
	AnswerID     string   `sql:"answer_id,pk" json:"answer_id"`
	ViewCount    int64    `sql:"view_count" json:"view_count"`
	LikeCount    int64    `sql:"like_count" json:"like_count"`
	DislikeCount int64    `sql:"dislike_count" json:"dislike_count"`
	LikeIndex    float64  `sql:"like_index" json:"like_index"`
}
