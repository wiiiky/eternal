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
