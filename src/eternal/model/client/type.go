package client

import (
	"time"
)

type Client struct {
	TableName   struct{}  `sql:"client" json:"-"`
	ID          string    `sql:"id,pk" json:"id"`
	Name        string    `sql:"name" json:"name"`
	TokenMaxAge uint64    `sql:"token_max_age" json:"token_max_age"`
	CTime       time.Time `sql:"ctime,null" json:"ctime"`
}
