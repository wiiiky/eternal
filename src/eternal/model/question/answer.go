package question

import (
	"eternal/errors"
	"eternal/model/db"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"time"
)

/* 获取问题下的回答 */
func GetQuestionAnswers(userID, questionID string, page, limit int) ([]*Answer, error) {
	conn := db.Conn()
	answers := make([]*Answer, 0)

	err := conn.Model(&answers).Column("User").Where("question_id = ?", questionID).Offset((page - 1) * limit).Limit(limit).Select()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return answers, nil
}

/* 获取热门回答 */
func FindHotAnswers(userID string, before string, limit int) ([]*HotAnswer, error) {
	conn := db.Conn()
	hotAnswers := make([]*HotAnswer, 0)

	var beforeTime time.Time
	var err error
	if before == "" {
		beforeTime = time.Now()
	} else {
		beforeTime, err = time.Parse(time.RFC3339, before)
		if err != nil {
			log.Warn("time.Parse failed:", err)
			return nil, err
		}
	}

	err = conn.Model(&hotAnswers).Column("hot_answer.*", "Answer", "Topic", "Question", "Answer.User").
		Where("hot_answer.ctime < ?", beforeTime).Limit(limit).Order("hot_answer.ctime DESC").Select()
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

/* 添加点赞，返回该回答的点赞数和踩数 */
func UpvoteAnswer(userID, answerID string) (uint64, uint64, error) {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}
	defer tx.Rollback()

	answer := Answer{
		ID: answerID,
	}
	if err := tx.Model(&answer).Column("id", "upvote_count", "downvote_count").WherePK().Select(); err != nil {
		if err == pg.ErrNoRows {
			return 0, 0, errors.ErrAnswerNotFound
		}
		log.Error("SQL Error", err)
		return 0, 0, err
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
		return answer.UpvoteCount, answer.DownvoteCount, nil
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	if err := tx.Select(&downvote); err == nil { /* 存在一个“踩“标签，删除它 */
		if err := tx.Delete(&downvote); err != nil {
			log.Error("SQL Error", err)
			return 0, 0, err
		}
		if _, err := tx.Model(&answer).Set("downvote_count = downvote_count - 1").WherePK().Returning("downvote_count").Update(); err != nil {
			log.Error("SQL Error", err)
			return 0, 0, err
		}
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	if err := tx.Insert(&upvote); err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	if _, err := tx.Model(&answer).Set("upvote_count = upvote_count + 1").WherePK().Returning("upvote_count").Update(); err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}
	return answer.UpvoteCount, answer.DownvoteCount, nil
}

/* 取消点赞 */
func UndoUpvoteAnswer(userID, answerID string) (uint64, uint64, error) {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}
	defer tx.Rollback()

	answer := Answer{
		ID: answerID,
	}
	if err := tx.Model(&answer).Column("id", "upvote_count", "downvote_count").WherePK().Select(); err != nil {
		if err == pg.ErrNoRows {
			return 0, 0, errors.ErrAnswerNotFound
		}
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	upvote := AnswerUpvote{
		UserID:   userID,
		AnswerID: answerID,
	}
	res, err := tx.Model(&upvote).WherePK().Delete()
	if err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}
	if res.RowsAffected() > 0 { /* 删除成功 */
		if _, err := tx.Model(&answer).Set("upvote_count = upvote_count - 1").WherePK().Returning("upvote_count").Update(); err != nil {
			log.Error("SQL Error", err)
			return 0, 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	return answer.UpvoteCount, answer.DownvoteCount, nil
}

/* 添加不喜欢 */
func DownvoteAnswer(userID, answerID string) (uint64, uint64, error) {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}
	defer tx.Rollback()

	answer := Answer{
		ID: answerID,
	}
	if err := tx.Model(&answer).Column("id", "upvote_count", "downvote_count").WherePK().Select(); err != nil {
		if err == pg.ErrNoRows {
			return 0, 0, errors.ErrAnswerNotFound
		}
		log.Error("SQL Error", err)
		return 0, 0, err
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
		return answer.UpvoteCount, answer.DownvoteCount, nil
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	if err := tx.Select(&upvote); err == nil { /* 存在一个“点赞“标签，删除“点赞“ */
		if err := tx.Delete(&upvote); err != nil {
			log.Error("SQL Error", err)
			return 0, 0, err
		}
		if _, err := tx.Model(&answer).Set("upvote_count = upvote_count - 1").WherePK().Returning("upvote_count").Update(); err != nil {
			log.Error("SQL Error", err)
			return 0, 0, err
		}
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	if err := tx.Insert(&downvote); err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	if _, err := tx.Model(&answer).Set("downvote_count = downvote_count + 1").WherePK().Returning("downvote_count").Update(); err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}
	return answer.UpvoteCount, answer.DownvoteCount, nil
}

/* 取消踩 */
func UndoDownvoteAnswer(userID, answerID string) (uint64, uint64, error) {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}
	defer tx.Rollback()

	answer := Answer{
		ID: answerID,
	}
	if err := tx.Model(&answer).Column("id", "upvote_count", "downvote_count").WherePK().Select(); err != nil {
		if err == pg.ErrNoRows {
			return 0, 0, errors.ErrAnswerNotFound
		}
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	downvote := AnswerDownvote{
		UserID:   userID,
		AnswerID: answerID,
	}
	res, err := tx.Model(&downvote).WherePK().Delete()
	if err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}
	if res.RowsAffected() > 0 { /* 删除成功 */
		if _, err := tx.Model(&answer).Set("downvote_count = downvote_count - 1").WherePK().Returning("downvote_count").Update(); err != nil {
			log.Error("SQL Error", err)
			return 0, 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Error("SQL Error", err)
		return 0, 0, err
	}

	return answer.UpvoteCount, answer.DownvoteCount, nil
}
