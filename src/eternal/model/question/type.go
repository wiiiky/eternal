package question

import (
	"eternal/model/account"
	"time"
)

type Topic struct {
	TableName   struct{}  `sql:"topic" json:"-"`
	ID          string    `sql:"id" json:"id"`
	Name        string    `sql:"name" json:"name"`
	Content     string    `sql:"content" json:"content"`
	UTime       time.Time `sql:"utime,null" json:"utime"`
	CTime       time.Time `sql:"ctime,null" json:"ctime"`
}

type Question struct {
	TableName   struct{}             `sql:"question" json:"-"`
	ID          string               `sql:"id", json:"id"`
	Title       string               `sql:"title" json:"title"`
	Content     string               `sql:"content" json:"content"`
	UserID      string               `sql:"user_id" json:"-"`
	UTime       time.Time            `sql:"utime,null" json:"utime"`
	CTime       time.Time            `sql:"ctime,null" json:"ctime"`
	Topics      []*Topic             `pg:"many2many:question_topic,fk:qid,joinFK:tid" json:"topics"`
	User        *account.UserProfile `sql:"-" json:"user"`
}

type QuestionTopic struct {
	TableName struct{}  `sql:"question_topic" json:"-"`
	QID       string    `sql:"qid,pk" json:"qid"`
	TID       string    `sql:"tid,pk" json:"tid"`
	CTime     time.Time `sql:"ctime,null" json:"ctime"`
}
