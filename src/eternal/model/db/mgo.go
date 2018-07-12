package db

import (
	"github.com/globalsign/mgo"
)

var _mgo *mgo.Session = nil
var _mdb *mgo.Database = nil

/* 初始化MongoDB */
func InitMongo(sURL, dbName string) error {
	var err error
	_mgo, err = mgo.Dial(sURL)
	if err != nil {
		return err
	}
	_mdb = _mgo.DB(dbName)
	return nil
}

func MGO() *mgo.Session {
	return _mgo
}

func MDB() *mgo.Database {
	return _mdb
}

/* 获取指定Collection */
func MC(c string) *mgo.Collection {
	return _mdb.C(c)
}
