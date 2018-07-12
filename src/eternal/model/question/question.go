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
		return nil, errors.ErrDB
	}
	return &question, nil
}

/* 获取用户和回答的关系信息 */
func GetUserQuestionRelationship(userID, questionID string) (*UserQuestionRelationship, error) {
	conn := db.Conn()

	userQuestionRelationship := &UserQuestionRelationship{
		Followed: false,
	}
	questionFollow := QuestionFollow{
		UserID:     userID,
		QuestionID: questionID,
	}
	if err := conn.Select(&questionFollow); err != nil {
		if err != pg.ErrNoRows {
			log.Error("SQL Error:", err)
			return nil, errors.ErrDB
		}
	} else {
		userQuestionRelationship.Followed = true
	}

	return userQuestionRelationship, nil
}

/* 搜索问题 */
func FindQuestions(query string, page, limit int) ([]*Question, error) {
	conn := db.Conn()

	questions := make([]*Question, 0)
	if err := conn.Model(&questions).Where("title LIKE ?", "%"+query+"%").Offset((page - 1) * limit).Limit(limit).Order("ctime DESC").Select(); err != nil {
		log.Error("SQL Error:", err)
		return nil, errors.ErrDB
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
		return nil, errors.ErrDB
	}
	for _, topicID := range topicIDs {
		topic := Topic{
			ID: topicID,
		}
		if err := tx.Model(&topic).Column("id").WherePK().Select(); err != nil {
			if err != pg.ErrNoRows {
				log.Error("SQL Error:", err)
				return nil, errors.ErrDB
			}
			return nil, errors.ErrTopicNotFound
		}
		qt := QuestionTopic{
			QuestionID: question.ID,
			TopicID:    topicID,
		}
		if err := tx.Insert(&qt); err != nil {
			log.Error("SQL Error:", err)
			return nil, errors.ErrDB
		}
	}
	if err := tx.Commit(); err != nil {
		log.Error("SQL Error:", err)
		return nil, errors.ErrDB
	}

	return GetQuestion(question.ID)
}

/* 关注问题 */
func FollowQuestion(userID, questionID string) (uint64, error) {
	tx, err := db.Conn().Begin()
	if err != nil {
		log.Error("SQL Error:", err)
		return 0, errors.ErrDB
	}
	defer tx.Rollback()

	question := Question{ID: questionID}
	if err := tx.Select(&question); err != nil {
		if err != pg.ErrNoRows {
			log.Error("SQL Error:", err)
			return 0, errors.ErrDB
		}
		return 0, errors.ErrQuestionNotFound
	}

	questionFollow := QuestionFollow{
		QuestionID: questionID,
		UserID:     userID,
	}
	if err := tx.Select(&questionFollow); err != nil {
		if err != pg.ErrNoRows {
			log.Error("SQL Error:", err)
			return 0, errors.ErrDB
		}
	} else {
		return question.FollowCount, nil
	}

	if err := tx.Insert(&questionFollow); err != nil {
		log.Error("SQL Error:", err)
		return 0, errors.ErrDB
	}

	if _, err := tx.Model(&question).Set("follow_count = follow_count + 1").Where("id = ?", questionID).Returning("follow_count").Update(); err != nil {
		log.Error("SQL Error:", err)
		return 0, errors.ErrDB
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Error:", err)
		return 0, errors.ErrDB
	}

	return question.FollowCount, nil
}

/* 取消问题的关注 */
func UnfollowQuestion(userID, questionID string) (uint64, error) {
	tx, err := db.Conn().Begin()
	if err != nil {
		log.Error("SQL Error:", err)
		return 0, errors.ErrDB
	}
	defer tx.Rollback()

	questionFollow := QuestionFollow{
		QuestionID: questionID,
		UserID:     userID,
	}
	if err := tx.Select(&questionFollow); err != nil {
		if err != pg.ErrNoRows {
			log.Error("SQL Error:", err)
			return 0, errors.ErrDB
		}
		return 0, nil
	}
	if err := tx.Delete(&questionFollow); err != nil {
		log.Error("SQL Error:", err)
		return 0, errors.ErrDB
	}

	question := Question{
		ID: questionID,
	}
	if _, err := tx.Model(&question).Set("follow_count = follow_count - 1").Where("id = ?", questionID).Returning("follow_count").Update(); err != nil {
		log.Error("SQL Error:", err)
		return 0, errors.ErrDB
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Error:", err)
		return 0, errors.ErrDB
	}
	return question.FollowCount, nil
}
