package file

import (
	"time"
)

type File struct {
	TableName   struct{}  `sql:"file" json:"-"`
	ID          string    `sql:"id" json:"id"`
	ContentType string    `sql:"content_type" json:"content_type"`
	CTime       time.Time `sql:"ctime,null" json:"ctime"`
}
