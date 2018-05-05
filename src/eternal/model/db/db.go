package db

import (
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
	"fmt"
	"time"
)

var _db *pg.DB = nil

func Conn() *pg.DB {
	return _db
}

func Start(sURL string) error {
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

			fmt.Printf("\033[34m%s %s\n\033[0m", time.Since(event.StartTime), query)
		})
	}
	return nil
}
