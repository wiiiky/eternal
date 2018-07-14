package account

import (
	"eternal/cache/store"
	accountModel "eternal/model/account"
	log "github.com/sirupsen/logrus"
	"time"
)

func GetToken(tokenID string) (*accountModel.Token, error) {
	var err error
	token := new(accountModel.Token)
	key := getTokenKey(tokenID)
	if ok, _ := store.GetVal(key, token); ok {
		return token, nil
	}
	if token, err = accountModel.GetToken(tokenID); err != nil {
		return nil, err
	}
	if err := store.SetVal(key, token, time.Hour * 24); err != nil {
		log.Error("SetVal failed:", err)
	}
	return token, err
}
