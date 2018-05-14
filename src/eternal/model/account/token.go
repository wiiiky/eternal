package account

import (
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
	"time"
)

func UpsertToken(userID string) (*Token, error) {
	conn := db.Conn()

	tk := &Token{
		UserID: userID,
		CTime:  time.Now(),
	}
	_, err := conn.Model(tk).
		OnConflict("(user_id) DO UPDATE").
		Set("id = uuid_generate_v1mc()").
		Set("ctime = ?ctime").Insert()
	if err != nil {
		log.Error("UpsertToken error:", err)
		return nil, err
	}
	return tk, nil
}

func DeleteToken(userID string) error {
	conn := db.Conn()
	tk := &Token{
		UserID: userID,
	}
	_, err := conn.Model(tk).Where("user_id = ?user_id").Delete()
	if err != nil {
		log.Error("DeleteToken  error:", err)
	}
	return err
}
