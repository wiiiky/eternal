package question

import (
	"eternal/model/db"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

func FindHotAnswers(userID string, page, limit int) ([]*HotAnswer, error) {
	conn := db.Conn()
	hotAnswers := make([]*HotAnswer, 0)

	err := conn.Model(&hotAnswers).Column("hot_answer.*", "Answer", "Topic", "Question", "Answer.User").
		Offset((page - 1) * limit).Limit(limit).Order("hot_answer.ctime DESC").Select()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	for _, hotAnswer := range hotAnswers {
		hotAnswer.Relationship = &UserAnswerRelationship{}
		upvote := AnswerUpvote{
			UserID:   userID,
			AnswerID: hotAnswer.Answer.ID,
		}
		if err := conn.Select(&upvote); err == nil {
			hotAnswer.Relationship.Upvoted = true
			continue // 如果存在点赞，则不继续查询是否有"踩"
		} else if err != pg.ErrNoRows {
			log.Error("SQL Error:", err)
			return nil, err
		}
		downvote := AnswerDownvote{
			UserID:   userID,
			AnswerID: hotAnswer.Answer.ID,
		}
		if err := conn.Select(&downvote); err == nil {
			hotAnswer.Relationship.Downvoted = true
		} else if err != pg.ErrNoRows {
			log.Error("SQL Error:", err)
			return nil, err
		}
	}

	return hotAnswers, nil
}

/* 添加喜欢 */
func UpvoteAnswer(userID, answerID string) error {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error", err)
		return err
	}
	defer tx.Rollback()

	answer := Answer{
		ID: answerID,
	}
	if err := tx.Model(&answer).Column("answer.id").Select(); err != nil {
		if err == pg.ErrNoRows {
			return db.ErrKeyNotFound
		}
		log.Error("SQL Error", err)
		return err
	}

	downvote := AnswerDownvote{
		UserID:   userID,
		AnswerID: answerID,
	}
	upvote := AnswerUpvote{
		UserID:   userID,
		AnswerID: answerID,
	}

	if err := tx.Select(&upvote); err == nil { /* “点赞“标签已经存在 */
		return nil
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Select(&downvote); err == nil { /* 存在一个“踩“标签，删除它 */
		if err := tx.Delete(&downvote); err != nil {
			log.Error("SQL Error", err)
			return err
		}
		if _, err := tx.Model(&answer).Set("downvote_count = downvote_count - 1").Where("id=?id").Update(); err != nil {
			log.Error("SQL Error", err)
			return err
		}
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Insert(&upvote); err != nil {
		log.Error("SQL Error", err)
		return err
	}

	if _, err := tx.Model(&answer).Set("upvote_count = upvote_count + 1").Where("id=?id").Update(); err != nil {
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Error", err)
		return err
	}
	return nil
}

/* 添加不喜欢 */
func DownvoteAnswer(userID, answerID string) error {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error", err)
		return err
	}
	defer tx.Rollback()

	answer := Answer{
		ID: answerID,
	}
	if err := tx.Model(&answer).Column("answer.id").Select(); err != nil {
		if err == pg.ErrNoRows {
			return db.ErrKeyNotFound
		}
		log.Error("SQL Error", err)
		return err
	}

	downvote := AnswerDownvote{
		UserID:   userID,
		AnswerID: answerID,
	}
	upvote := AnswerUpvote{
		UserID:   userID,
		AnswerID: answerID,
	}

	if err := tx.Select(&downvote); err == nil { /* “踩“已经存在 */
		return nil
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Select(&upvote); err == nil { /* 存在一个“点赞“标签，删除“点赞“ */
		if err := tx.Delete(&upvote); err != nil {
			log.Error("SQL Error", err)
			return err
		}
		if _, err := tx.Model(&answer).Set("upvote_count = upvote_count - 1").Where("id=?id").Update(); err != nil {
			log.Error("SQL Error", err)
			return err
		}
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Insert(&downvote); err != nil {
		log.Error("SQL Error", err)
		return err
	}

	if _, err := tx.Model(&answer).Set("downvote_count = downvote_count + 1").Where("id=?id").Update(); err != nil {
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Error", err)
		return err
	}
	return nil
}
