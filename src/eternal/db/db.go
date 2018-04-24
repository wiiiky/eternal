package db

import (
	"github.com/go-pg/pg"
)

var _db *pg.DB = nil

func DB() *pg.DB {
	return _db
}

func Start(sURL string) error {
	opts, err := pg.ParseURL(sURL)
	if err != nil {
		return err
	}
	_db = pg.Connect(opts)
	return nil
}
