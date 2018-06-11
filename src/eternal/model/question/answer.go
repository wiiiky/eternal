package question

import (
	"eternal/errors"
	"eternal/model/db"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"time"
)

func GetAnswer(answerID string) (*Answer, error) {
	conn := db.Conn()
	answer := Answer{
		ID: answerID,
	}
	if err := conn.Select(&answer); err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		log.Error("SQL Error:", err)
		return nil, err
	}
	return &answer, nil
}

/* 获取用户和某个回答的关系 */
func GetUserAnswerRelationship(userID, answerID string) (*UserAnswerRelationship, error) {
	conn := db.Conn()
	relationship := &UserAnswerRelationship{}
	upvote := AnswerUpvote{
		UserID:   userID,
		AnswerID: answerID,
	}
	if err := conn.Select(&upvote); err == nil {
		relationship.Upvoted = true
		return relationship, nil
	} else if err != pg.ErrNoRows {
		log.Error("SQL Error:", err)
		return nil, err
	}
	downvote := AnswerDownvote{
		UserID:   userID,
		AnswerID: answerID,
	}
	if err := conn.Select(&downvote); err == nil {
		relationship.Downvoted = true
	} else if err != pg.ErrNoRows {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return relationship, nil
}

/* 获取问题下的回答 */
func GetQuestionAnswers(userID, questionID string, page, limit int) ([]*Answer, error) {
	conn := db.Conn()
	answers := make([]*Answer, 0)

	err := conn.Model(&answers).Column("User").Where("question_id = ?", questionID).Order("upvote_count DESC").Offset((page - 1) * limit).Limit(limit).Select()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return answers, nil
}

/* 获取问题的最热门回答 */
func GetQuestionTopAnswer(questionID string) (*Answer, error) {
	conn := db.Conn()
	answer := Answer{}
	err := conn.Model(&answer).Where("question_id = ?", questionID).Order("upvote_count DESC").Limit(1).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return &answer, nil
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

	return hotAnswers, nil
}

/* 添加点赞，返回该回答的点赞数和踩数 */
func UpvoteAnswer(userID, answerID string) (uint64, uint64, error) {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error:", err)
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
		log.Error("SQL Error:", err)
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
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	if err := tx.Select(&downvote); err == nil { /* 存在一个“踩“标签，删除它 */
		if err := tx.Delete(&downvote); err != nil {
			log.Error("SQL Error:", err)
			return 0, 0, err
		}
		if _, err := tx.Model(&answer).Set("downvote_count = downvote_count - 1").WherePK().Returning("downvote_count").Update(); err != nil {
			log.Error("SQL Error:", err)
			return 0, 0, err
		}
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	if err := tx.Insert(&upvote); err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	if _, err := tx.Model(&answer).Set("upvote_count = upvote_count + 1").WherePK().Returning("upvote_count").Update(); err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}
	return answer.UpvoteCount, answer.DownvoteCount, nil
}

/* 取消点赞 */
func UndoUpvoteAnswer(userID, answerID string) (uint64, uint64, error) {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error:", err)
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
		log.Error("SQL Error:", err)
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
			log.Error("SQL Error:", err)
			return 0, 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	return answer.UpvoteCount, answer.DownvoteCount, nil
}

/* 添加不喜欢 */
func DownvoteAnswer(userID, answerID string) (uint64, uint64, error) {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error:", err)
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
		log.Error("SQL Error:", err)
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
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	if err := tx.Select(&upvote); err == nil { /* 存在一个“点赞“标签，删除“点赞“ */
		if err := tx.Delete(&upvote); err != nil {
			log.Error("SQL Error:", err)
			return 0, 0, err
		}
		if _, err := tx.Model(&answer).Set("upvote_count = upvote_count - 1").WherePK().Returning("upvote_count").Update(); err != nil {
			log.Error("SQL Error:", err)
			return 0, 0, err
		}
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	if err := tx.Insert(&downvote); err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	if _, err := tx.Model(&answer).Set("downvote_count = downvote_count + 1").WherePK().Returning("downvote_count").Update(); err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}
	return answer.UpvoteCount, answer.DownvoteCount, nil
}

/* 取消踩 */
func UndoDownvoteAnswer(userID, answerID string) (uint64, uint64, error) {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error:", err)
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
		log.Error("SQL Error:", err)
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
			log.Error("SQL Error:", err)
			return 0, 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Error("SQL Error:", err)
		return 0, 0, err
	}

	return answer.UpvoteCount, answer.DownvoteCount, nil
}

/* 获取在指定时间内点赞次数 */
func GetAnswerUpvoteCount(answerID string, startTime, endTime time.Time) (int, error) {
	conn := db.Conn()
	count, err := conn.Model((*AnswerUpvote)(nil)).Where("answer_id = ? AND ctime > ? AND ctime < ?", answerID, startTime, endTime).CountEstimate(1000)
	if err != nil {
		log.Error("SQL Error:", err)
	}
	return count, err
}

/* 获取在指定时间内的踩次数 */
func GetAnswerDownvoteCount(answerID string, startTime, endTime time.Time) (int, error) {
	conn := db.Conn()
	count, err := conn.Model((*AnswerDownvote)(nil)).Where("answer_id = ? AND ctime > ? AND ctime < ?", answerID, startTime, endTime).CountEstimate(1000)
	if err != nil {
		log.Error("SQL Error:", err)
	}
	return count, err
}

/* 添加或者更新热门回答 */
func UpsertHotAnswer(answerID string) error {
	conn := db.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Error("SQL Error:", err)
		return err
	}
	defer tx.Rollback()

	answer := Answer{
		ID: answerID,
	}
	if err := tx.Model(&answer).Column("Question", "Question.Topics").WherePK().Select(); err != nil {
		if err != pg.ErrNoRows {
			log.Error("SQL Error:", err)
		}
		return errors.ErrAnswerNotFound
	}
	/* 删除同一个问题下其他的热门回答，同一个问题只能有一个热门回答 */
	if _, err := tx.Model((*HotAnswer)(nil)).Where("question_id = ?", answer.Question.ID).Delete(); err != nil {
		log.Error("SQL Error:", err)
		return err
	}
	for _, topic := range answer.Question.Topics {
		hotAnswer := HotAnswer{}
		err := tx.Model(&hotAnswer).Where("answer_id = ? AND question_id = ? AND topic_id = ?", answer.ID, answer.Question.ID, topic.ID).Select()
		if err == nil {
			_, err = tx.Model(&hotAnswer).Set("ctime = ?", time.Now()).WherePK().Update()
			if err != nil {
				log.Error("SQL Error:", err)
				return err
			}
		} else if err == pg.ErrNoRows {
			hotAnswer.AnswerID = answer.ID
			hotAnswer.QuestionID = answer.Question.ID
			hotAnswer.TopicID = topic.ID
			if err := tx.Insert(&hotAnswer); err != nil {
				log.Error("SQL Error:", err)
				return err
			}
		} else {
			log.Error("SQL Error:", err)
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Error("SQL Error:", err)
		return err
	}
	return nil
}

/* 删除热门回答 */
func DeleteHotAnswer(answerID string) error {
	conn := db.Conn()

	if _, err := conn.Model((*HotAnswer)(nil)).Where("answer_id = ?", answerID).Delete(); err != nil {
		log.Error("SQL Error:", err)
		return err
	}
	return nil
}
