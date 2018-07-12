package client

import (
	"eternal/cache/store"
	clientModel "eternal/model/client"
	log "github.com/sirupsen/logrus"
)

func GetClient(clientID string) (*clientModel.Client, error) {
	var err error
	client := new(clientModel.Client)
	if ok, _ := store.HGetVal(KeyClient, clientID, client); ok {
		return client, nil
	}
	if client, err = clientModel.GetClient(clientID); err != nil {
		return nil, err
	}
	if err := store.HSetVal(KeyClient, client.ID, client); err != nil {
		log.Error("HSet failed:", err)
	}
	return client, err
}
