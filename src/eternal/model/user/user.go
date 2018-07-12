package user

import (
	"eternal/model/account"
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
)

func GetUserProfile(userID string) (*account.UserProfile, error) {
	conn := db.PG()

	up := &account.UserProfile{
		UserID: userID,
	}
	err := conn.Select(up)
	if err == db.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return up, nil
}

func UpdateUserCover(userID, cover string) (*account.UserProfile, error) {
	conn := db.PG()

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
