package account

import (
	"time"
)

const (
	PTYPE_MD5    = "MD5"
	PTYPE_SHA1   = "SHA1"
	PTYPE_SHA256 = "SHA256"
)

type SupportedCounty struct {
	TableName struct{} `sql:"supported_country" json:"-"`
	Code      string   `sql:"code,pk" json:"code"`
	Name      string   `sql:"name" json:"name"`
}

type Account struct {
	TableName   struct{}  `sql:"account" json:"-"`
	ID          string    `sql:"id" json:"id"`
	CountryCode string    `sql:"country_code" json:"country_code"`
	PhoneNumber string    `sql:"phone_number" json:"phone_number"`
	Salt        string    `sql:"salt" json:"-"`
	Password    string    `sql:"passwd" json:"-"`
	PType       string    `sql:"ptype" json:"-"`
	UTime       time.Time `sql:"utime,null" json:"utime"`
	CTime       time.Time `sql:"ctime,null" json:"ctime"`
}

func (a *Account) Auth(plain string) bool {
	password, _ := encryptPassword(a.Salt, plain, a.PType)
	return password == a.Password
}

type Token struct {
	TableName struct{}  `sql:"token" json:"-"`
	ID        string    `sql:"id" json:"id"`
	UserID    string    `sql:"user_id" json:"user_id"`
	ClientID  string    `sql:"client_id" json:"client_id"`
	CTime     time.Time `sql:"ctime" json:"ctime"`
}

type UserProfile struct {
	TableName   struct{}  `sql:"user_profile" json:"-"`
	UserID      string    `sql:"user_id,pk" json:"user_id"`
	Name        string    `sql:"name" json:"name"`
	Gender      string    `sql:"gender" json:"gender"`
	Description string    `sql:"description" json:"description"`
	Avatar      string    `sql:"avatar" json:"avatar"`
	Cover       string    `sql:"cover" json:"cover"`
	Birthday    time.Time `sql:"birthday" json:"birthday"`
	UTime       time.Time `sql:"utime,null" json:"utime"`
	CTime       time.Time `sql:"ctime,null" json:"ctime"`
}
