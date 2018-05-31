package question

import (
	"eternal/errors"
	"eternal/model/db"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

func GetQuestion(id string) (*Question, error) {
	conn := db.Conn()

	question := Question{
		ID: id,
	}
	err := conn.Model(&question).Column("question.*", "User", "Topics").WherePK().Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return &question, nil
}

/* 搜索问题 */
func FindQuestions(query string, page, limit int) ([]*Question, error) {
	conn := db.Conn()

	questions := make([]*Question, 0)
	if err := conn.Model(&questions).Where("title LIKE ?", "%"+query+"%").Offset((page - 1) * limit).Limit(limit).Order("ctime DESC").Select(); err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return questions, nil
}

/* 创建问题 */
func CreateQuestion(userID, title string, topicIDs []string, content string) (*Question, error) {
	tx, err := db.Conn().Begin()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	defer tx.Rollback()

	question := Question{
		Title:   title,
		Content: content,
		UserID:  userID,
	}
	if err := tx.Insert(&question); err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	for _, topicID := range topicIDs {
		topic := Topic{
			ID: topicID,
		}
		if err := tx.Model(&topic).Column("id").WherePK().Select(); err != nil {
			if err != pg.ErrNoRows {
				log.Error("SQL Error:", err)
			} else {
				log.Warnf("Topic %s Not Found.", topicID)
			}
			return nil, errors.ErrTopicNotFound
		}
		qt := QuestionTopic{
			QuestionID: question.ID,
			TopicID:    topicID,
		}
		if err := tx.Insert(&qt); err != nil {
			log.Error("SQL Error:", err)
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}

	return GetQuestion(question.ID)
}
