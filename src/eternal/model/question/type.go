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
	Topics    []*Topic                  `pg:"many2many:question_topic,fk:question_id,joinFK:topic_id" json:"topics"`
	User      *accountModel.UserProfile `sql:"-" json:"user"`
}

type QuestionTopic struct {
	TableName   struct{}  `sql:"question_topic" json:"-"`
	QuesetionID string    `sql:"question_id,pk" json:"question_id"`
	TopicID     string    `sql:"topic_id,pk" json:"topic_id"`
	CTime       time.Time `sql:"ctime,null" json:"ctime"`
}

type Answer struct {
	TableName     struct{}                  `sql:"answer" json:"-"`
	ID            string                    `sql:"id" json:"id"`
	Content       string                    `sql:"content" json:"content"`
	Excerpt       string                    `sql:"excerpt" json:"excerpt"`
	QuestionID    string                    `sql:"question_id" json:"question_id"`
	UserID        string                    `sql:"user_id" json:"-"`
	ViewCount     int64                     `sql:"view_count" json:"view_count"`
	UpvoteCount   int64                     `sql:"upvote_count" json:"upvote_count"`
	DownvoteCount int64                     `sql:"downvote_count" json:"downvote_count"`
	UTime         time.Time                 `sql:"utime,null" json:"utime"`
	CTime         time.Time                 `sql:"ctime,null" json:"ctime"`
	User          *accountModel.UserProfile `json:"user"`
	Question      *Question                 `json:"question"`
}

/* 包含用户和回答的关系信息，不对应任何具体数据表 */
type UserAnswerRelationship struct {
	Upvoted   bool `json:"upvoted"`
	Downvoted bool `json:"downvoted"`
}

type HotAnswer struct {
	TableName    struct{}                `sql:"hot_answer" json:"-"`
	ID           string                  `sql:"id,pk" json:"-"`
	AnswerID     string                  `sql:"answer_id" json:"-"`
	QuestionID   string                  `sql:"question_id" json:"-"`
	TopicID      string                  `sql:"topic_id" json:"-"`
	Answer       *Answer                 `json:"answer"`
	Question     *Question               `json:"question"`
	Topic        *Topic                  `json:"topic"`
	CTime        time.Time               `sql:"ctime,null" json:"-"`
	Relationship *UserAnswerRelationship `json:"relationship"`
}

type AnswerUpvote struct {
	TableName struct{}  `sql:"answer_upvote" json:"-"`
	UserID    string    `sql:"user_id,pk" json:"answer_id"`
	AnswerID  string    `sql:"answer_id,pk" json:"answer_id"`
	CTime     time.Time `sql:"ctime,null" json:"ctime"`
}

type AnswerDownvote struct {
	TableName struct{}  `sql:"answer_downvote" json:"-"`
	UserID    string    `sql:"user_id,pk" json:"answer_id"`
	AnswerID  string    `sql:"answer_id,pk" json:"answer_id"`
	CTime     time.Time `sql:"ctime,null" json:"ctime"`
}
