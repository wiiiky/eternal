package client

import (
	"eternal/cache/store"
	clientModel "eternal/model/client"
	log "github.com/sirupsen/logrus"
)

func GetClient(clientID string) (*clientModel.Client, error) {
	var err error
	client := new(clientModel.Client)

	key := getClientKey(clientID)
	if ok, _ := store.GetVal(key, client); ok {
		return client, nil
	}
	if client, err = clientModel.GetClient(clientID); err != nil {
		return nil, err
	}
	if err = store.SetVal(key, client, 0); err != nil {
		log.Error("HSet failed:", err)
	}
	return client, nil
}
