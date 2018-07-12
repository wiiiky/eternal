package client

import (
	"eternal/errors"
	"eternal/model/db"
	log "github.com/sirupsen/logrus"
)

func GetClient(clientID string) (*Client, error) {
	conn := db.PG()

	client := &Client{
		ID: clientID,
	}
	if err := conn.Select(client); err != nil {
		if err == db.ErrNoRows {
			return nil, nil
		}
		log.Error("SQL Error:", err)
		return nil, errors.ErrDB
	}
	return client, nil
}
