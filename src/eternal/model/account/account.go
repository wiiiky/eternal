package account

import (
	"eternal/model/db"
	"eternal/util"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"time"
)

func GetSupportedCountries() ([]*SupportedCounty, error) {
	conn := db.Conn()

	countries := make([]*SupportedCounty, 0)
	err := conn.Model(&countries).Order(`sort ASC`).Select()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return countries, nil
}

func GetSupportedCountryWithCode(code string) (*SupportedCounty, error) {
	conn := db.Conn()
	country := &SupportedCounty{
		Code: code,
	}

	err := conn.Select(country)
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return country, nil
}

/* 创建帐号 */
func CreateAccount(countryCode, mobile, password, ptype string) (*Account, error) {
	conn := db.Conn()

	tx, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	a := &Account{
		CountryCode: countryCode,
		Mobile:      mobile,
	}
	err = tx.Model(a).Where(`mobile=?`, mobile).Select()
	if err != nil {
		if err != pg.ErrNoRows {
			return nil, err
		}
	} else { /* 查询成功，手机号已存在 */
		return nil, db.ErrKeyDuplicate
	}

	salt, err := util.RandString(12)
	if err != nil {
		log.Error("RandString error:", err)
		return nil, err
	}

	password, ptype = encryptPassword(salt, password, ptype)
	a.Salt = salt
	a.Password = password
	a.PType = ptype
	if err := tx.Insert(a); err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}

	up := &UserProfile{
		UserID: a.ID,
	}
	if err := tx.Insert(up); err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Error("SQL Commit failed:", err)
		return nil, err
	}
	log.Infof("User %s:%s created", countryCode, mobile)
	return a, nil
}

func GetAccountWithMobile(mobile string) (*Account, error) {
	conn := db.Conn()

	a := &Account{}
	err := conn.Model(a).Where("mobile = ?", mobile).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error: ", err)
		return nil, err
	}
	return a, err
}

func GetAccountWithUserID(userID string) (*Account, error) {
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

func GetAccountWithTokenID(tokenID string) (*Account, error) {
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
