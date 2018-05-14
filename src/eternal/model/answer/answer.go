package answer

import (
	"eternal/model/db"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

/* 添加喜欢 */
func AddAnswerLike(userID, answerID string) error {
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

	dislike := AnswerDislike{
		UserID:   userID,
		AnswerID: answerID,
	}
	like := AnswerLike{
		UserID:   userID,
		AnswerID: answerID,
	}
	answerStat := AnswerStat{
		AnswerID: answerID,
	}

	if err := tx.Select(&like); err == nil { /* “喜欢“标签已经存在 */
		return nil
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Select(&dislike); err == nil { /* 存在一个“不喜欢“标签，删除不喜欢 */
		if err := tx.Delete(&dislike); err != nil {
			log.Error("SQL Error", err)
			return err
		}
		if _, err := tx.Model(&answerStat).Set("dislike_count = dislike_count - 1").Where("answer_id=?answer_id").Update(); err != nil {
			log.Error("SQL Error", err)
			return err
		}
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Insert(&like); err != nil {
		log.Error("SQL Error", err)
		return err
	}

	if _, err := tx.Model(&answerStat).Set("like_count = like_count + 1").Where("answer_id=?answer_id").Update(); err != nil {
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
func AddAnswerDislike(userID, answerID string) error {
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

	dislike := AnswerDislike{
		UserID:   userID,
		AnswerID: answerID,
	}
	like := AnswerLike{
		UserID:   userID,
		AnswerID: answerID,
	}
	answerStat := AnswerStat{
		AnswerID: answerID,
	}

	if err := tx.Select(&dislike); err == nil { /* “不喜欢“标签已经存在 */
		return nil
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Select(&like); err == nil { /* 存在一个“喜欢“标签，删除喜欢 */
		if err := tx.Delete(&like); err != nil {
			log.Error("SQL Error", err)
			return err
		}
		if _, err := tx.Model(&answerStat).Set("like_count = like_count - 1").Where("answer_id=?answer_id").Update(); err != nil {
			log.Error("SQL Error", err)
			return err
		}
	} else if err != pg.ErrNoRows { /* 出错 */
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Insert(&dislike); err != nil {
		log.Error("SQL Error", err)
		return err
	}

	if _, err := tx.Model(&answerStat).Set("dislike_count = dislike_count + 1").Where("answer_id=?answer_id").Update(); err != nil {
		log.Error("SQL Error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Error", err)
		return err
	}
	return nil
}
