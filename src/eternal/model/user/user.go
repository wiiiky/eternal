package user

import (
	"eternal/model/account"
	"eternal/model/db"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

func GetUserProfile(userID string) (*account.UserProfile, error) {
	conn := db.Conn()

	up := &account.UserProfile{
		UserID: userID,
	}
	err := conn.Select(up)
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return up, nil
}

func UpdateUserCover(userID, cover string) (*account.UserProfile, error) {
	conn := db.Conn()

	up := &account.UserProfile{
		UserID: userID,
		Cover:  cover,
	}
	_, err := conn.Model(up).Column("cover").WherePK().Update()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return GetUserProfile(userID)
}
