package client

import (
	"eternal/errors"
	"eternal/model/db"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
)

func GetClient(clientID string) (*Client, error) {
	conn := db.Conn()

	client := &Client{
		ID: clientID,
	}
	if err := conn.Select(client); err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		log.Error("SQL Error:", err)
		return nil, errors.ErrDB
	}
	return client, nil
}
