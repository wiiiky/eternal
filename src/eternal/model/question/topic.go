package question

import (
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
)

func FindTopics(query string, limit int) ([]*Topic, error) {
	conn := db.Conn()

	topics := make([]*Topic, 0)
	err := conn.Model(&topics).Where("name LIKE ?", "%"+query+"%").Limit(limit).Select()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return topics, nil
}
