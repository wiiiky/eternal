package file

import (
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
)

/* 保存文件 */
func SaveFile(pk, contentType string) (*File, error) {
	conn := db.PG()

	f := &File{
		ID:          pk,
		ContentType: contentType,
	}

	if _, err := conn.Model(f).OnConflict("(id) DO UPDATE").Set("content_type = ?content_type").Insert(); err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return f, nil
}

func GetFile(pk string) (*File, error) {
	conn := db.PG()

	f := &File{
		ID: pk,
	}

	err := conn.Select(f)
	if err == db.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error("SQL Error:", err)
		return nil, err
	}
	return f, nil
}
