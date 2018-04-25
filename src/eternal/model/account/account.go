package account

import (
	"eternal/model/db"
	"eternal/util"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"time"
)

func Create(countryCode, mobile, password, ptype string) (*Account, error) {
	conn := db.Conn()

	salt, err := util.RandString(12)
	if err != nil {
		log.Error("RandString error:", err)
		return nil, err
	}

	password, ptype = encryptPassword(salt, password, ptype)

	a := &Account{
		CountryCode: countryCode,
		Mobile:      mobile,
		Salt:        salt,
		Password:    password,
		PType:       ptype,
	}
	if err := conn.Insert(a); err != nil {
		log.Error("SQL Error: ", err)
		return nil, err
	}
	log.Infof("User %s:%s created", countryCode, mobile)
	return a, nil
}

func GetWithMobile(countryCode, mobile string) (*Account, error) {
	conn := db.Conn()

	a := &Account{}
	err := conn.Model(a).Where("country_code = ?", countryCode).Where("mobile = ?", mobile).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error: ", err)
		return nil, err
	}
	return a, err
}

func GetWithUserID(userID string) (*Account, error) {
	conn := db.Conn()

	a := &Account{}
	err := conn.Model(a).Where("id = ?", userID).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error: ", err)
		return nil, err
	}
	return a, err
}

func GetWithTokenID(tokenID string) (*Account, error) {
	conn := db.Conn()

	a := &Account{}
	err := conn.Model(a).Join("JOIN token ON token.user_id=account.id").Where("token.id = ?", tokenID).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error: ", err)
		return nil, err
	}
	return a, nil
}

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
