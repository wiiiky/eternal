package question

import (
	accountModel "eternal/model/account"
	"time"
)

type Topic struct {
	TableName    struct{}  `sql:"topic" json:"-"`
	ID           string    `sql:"id" json:"id"`
	Name         string    `sql:"name" json:"name"`
	Introduction string    `sql:"introduction" json:"introduction"`
	UTime        time.Time `sql:"utime,null" json:"utime"`
	CTime        time.Time `sql:"ctime,null" json:"ctime"`
}

type Question struct {
	TableName struct{}                  `sql:"question" json:"-"`
	ID        string                    `sql:"id", json:"id"`
	Title     string                    `sql:"title" json:"title"`
	Content   string                    `sql:"content" json:"content"`
	UserID    string                    `sql:"user_id" json:"-"`
	UTime     time.Time                 `sql:"utime,null" json:"utime"`
	CTime     time.Time                 `sql:"ctime,null" json:"ctime"`
	Topics    []*Topic                  `pg:"many2many:question_topic,fk:qid,joinFK:tid" json:"topics"`
	User      *accountModel.UserProfile `sql:"-" json:"user"`
}

type QuestionTopic struct {
	TableName struct{}  `sql:"question_topic" json:"-"`
	QID       string    `sql:"qid,pk" json:"qid"`
	TID       string    `sql:"tid,pk" json:"tid"`
	CTime     time.Time `sql:"ctime,null" json:"ctime"`
}

type Answer struct {
	TableName    struct{}                  `sql:"answer" json:"-"`
	ID           string                    `sql:"id" json:"id"`
	Content      string                    `sql:"content" json:"content"`
	QuestionID   string                    `sql:"question_id" json:"question_id"`
	UserID       string                    `sql:"user_id" json:"-"`
	ViewCount    int64                     `sql:"view_count" json:"view_count"`
	LikeCount    int64                     `sql:"like_count" json:"like_count"`
	DislikeCount int64                     `sql:"dislike_count" json:"dislike_count"`
	UTime        time.Time                 `sql:"utime,null" json:"utime"`
	CTime        time.Time                 `sql:"ctime,null" json:"ctime"`
	User         *accountModel.UserProfile `json:"user"`
	Question     *Question                 `json:"question"`
}

type HotAnswer struct {
	TableName  struct{}  `sql:"hot_answer" json:"-"`
	ID         string    `sql:"id,pk" json:"-"`
	AnswerID   string    `sql:"answer_id" json:"-"`
	QuestionID string    `sql:"question_id" json:"-"`
	TopicID    string    `sql:"topic_id" json:"-"`
	Answer     *Answer   `json:"answer"`
	Question   *Question `json:"question"`
	Topic      *Topic    `json:"topic"`
	CTime      time.Time `sql:"ctime,null" json:"-"`
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
