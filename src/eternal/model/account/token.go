package account

import (
	"eternal/errors"
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
	"time"
)

func GetToken(tokenID string) (*Token, error) {
	conn := db.PG()

	tk := Token{
		ID: tokenID,
	}
	if err := conn.Select(&tk); err != nil {
		if err != db.ErrNoRows {
			log.Error("SQL Error:", err)
			return nil, errors.ErrDB
		}
		return nil, nil
	}
	return &tk, nil
}

func UpsertToken(userID, clientID string, expire time.Duration) (*Token, error) {
	conn := db.PG()

	now := time.Now()
	etime := now.Add(expire)
	tk := Token{
		UserID:   userID,
		ClientID: clientID,
		ETime:    etime,
		CTime:    now,
	}
	_, err := conn.Model(&tk).
		OnConflict("(user_id,client_id) DO UPDATE").
		Set("id = uuid_generate_v1mc()").
		Set("etime = ?etime").
		Set("ctime = ?ctime").Insert()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, errors.ErrDB
	}
	return &tk, nil
}

func DeleteToken(tokenID string) error {
	conn := db.PG()
	_, err := conn.Model((*Token)(nil)).Where("id = ?", tokenID).Delete()
	if err != nil {
		log.Error("SQL Error:", err)
		return errors.ErrDB
	}
	return err
}
