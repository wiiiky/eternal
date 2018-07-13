package account

import (
	"eternal/errors"
	"eternal/model/db"
	"eternal/util"
	log "github.com/sirupsen/logrus"
)

func GetSupportedCountries() ([]*SupportedCounty, error) {
	conn := db.PG()

	countries := make([]*SupportedCounty, 0)
	err := conn.Model(&countries).Order(`sort ASC`).Select()
	if err != nil {
		log.Error("SQL Error:", err)
		return nil, errors.ErrDB
	}
	return countries, nil
}

func GetSupportedCountryWithCode(code string) (*SupportedCounty, error) {
	conn := db.PG()
	country := &SupportedCounty{
		Code: code,
	}

	err := conn.Select(country)
	if err == db.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error:", err)
		return nil, errors.ErrDB
	}
	return country, nil
}

/* 创建帐号 */
func CreateAccount(countryCode, phoneNumber, password, ptype string) (*Account, error) {
	conn := db.PG()

	tx, err := conn.Begin()
	if err != nil {
		return nil, errors.ErrDB
	}
	defer tx.Rollback()

	a := &Account{
		CountryCode: countryCode,
		PhoneNumber: phoneNumber,
	}
	if err := tx.Model(a).Where(`phone_number=?`, phoneNumber).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Error("SQL Error:", err)
			return nil, errors.ErrDB
		}
	} else { /* 查询成功，手机号已存在 */
		return nil, errors.ErrPhoneNumberExisted
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
	log.Infof("User %s:%s created", countryCode, phoneNumber)
	return a, nil
}

func GetAccountWithPhoneNumber(phoneNumber string) (*Account, error) {
	conn := db.PG()

	a := &Account{}
	err := conn.Model(a).Where("phone_number = ?", phoneNumber).Select()
	if err == db.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error: ", err)
		return nil, err
	}
	return a, err
}

func GetAccountWithUserID(userID string) (*Account, error) {
	conn := db.PG()

	a := &Account{}
	err := conn.Model(a).Where("id = ?", userID).Select()
	if err == db.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error: ", err)
		return nil, err
	}
	return a, err
}

func GetAccountWithTokenID(tokenID string) (*Account, error) {
	conn := db.PG()

	a := &Account{}
	err := conn.Model(a).Join("JOIN token ON token.user_id=account.id").Where("token.id = ?", tokenID).Select()
	if err == db.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error: ", err)
		return nil, err
	}
	return a, nil
}
