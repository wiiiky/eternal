package account

import (
	"eternal/cache/store"
	accountModel "eternal/model/account"
	log "github.com/sirupsen/logrus"
)

func GetAccount(userID string) (*accountModel.Account, error) {
	var err error
	account := new(accountModel.Account)

	key := getAccountKey(userID)
	if ok, _ := store.GetVal(key, account); ok {
		return account, nil
	}
	if account, err = accountModel.GetAccount(userID); err != nil {
		return nil, err
	}
	if err = store.SetVal(key, account, 0); err != nil {
		log.Error("HSet failed:", err)
	}
	return account, nil
}
