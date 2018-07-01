package account

import (
	"eternal/errors"
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
	"time"
)

func UpsertToken(userID, clientID string) (*Token, error) {
	conn := db.Conn()

	tk := &Token{
		UserID:   userID,
		ClientID: clientID,
		CTime:    time.Now(),
	}
	_, err := conn.Model(tk).
		OnConflict("(user_id,client_id) DO UPDATE").
		Set("id = uuid_generate_v1mc()").
		Set("ctime = ?ctime").Insert()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, errors.ErrDB
	}
	return tk, nil
}

func DeleteToken(userID, clientID string) error {
	conn := db.Conn()
	_, err := conn.Model((*Token)(nil)).Where("user_id = ?", userID).Where("client_id = ?", clientID).Delete()
	if err != nil {
		log.Error("SQL Error:", err)
		return errors.ErrDB
	}
	return err
}
