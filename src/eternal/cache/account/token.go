package account

import (
	"eternal/cache/store"
	accountModel "eternal/model/account"
	log "github.com/sirupsen/logrus"
)

func GetToken(tokenID string) (*accountModel.Token, error) {
	var err error
	token := new(accountModel.Token)
	if ok, _ := store.HGetVal(KeyToken, tokenID, token); ok {
		return token, nil
	}
	if token, err = accountModel.GetToken(tokenID); err != nil {
		return nil, err
	}
	if err := store.HSetVal(KeyToken, token.ID, token); err != nil {
		log.Error("HSet failed:", err)
	}
	return token, err
}
