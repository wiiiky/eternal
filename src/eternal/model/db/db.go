package db

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
	"time"
)

var _db *pg.DB = nil

func Conn() *pg.DB {
	return _db
}

func Init(sURL string) error {
	opts, err := pg.ParseURL(sURL)
	if err != nil {
		return err
	}
	_db = pg.Connect(opts)
	/* DEBUG */
	debug := viper.GetBool("debug")
	if debug {
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
