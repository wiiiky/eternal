package question

import (
	"eternal/model/db"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	log "github.com/sirupsen/logrus"
)

func FindQuestions(userID string) ([]*Question, error) {
	conn := db.Conn()

	questions := make([]*Question, 0)
	err := conn.Model(&questions).Column("User").Column("question.*", "Topics").Relation("Topics", func(q *orm.Query) (*orm.Query, error) {
		return q.Order("qid"), nil
	}).OrderExpr("utime DESC").Limit(10).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return questions, nil
}
