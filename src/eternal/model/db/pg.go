package db

import (
	"eternal/config"
	"fmt"
	"github.com/go-pg/pg"
	"time"
)

var _db *pg.DB = nil

var ErrNoRows = pg.ErrNoRows

func PG() *pg.DB {
	return _db
}

/* 初始化PostgreSQL */
func InitPG(sURL string) error {
	opts, err := pg.ParseURL(sURL)
	if err != nil {
		return err
	}
	_db = pg.Connect(opts)
	/* DEBUG */
	if config.DEBUG {
		_db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}
			/* 打印SQL执行事件 */
			fmt.Printf("\033[34m%s %s\n\033[0m", time.Since(event.StartTime), query)
		})
	}
	return nil
}
