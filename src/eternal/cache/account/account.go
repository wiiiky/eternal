package account

import (
	"eternal/cache/store"
	accountModel "eternal/model/account"
	log "github.com/sirupsen/logrus"
)

func GetAccount(accountID string) (*accountModel.Account, error) {
	var err error
	account := new(accountModel.Account)
	if ok, _ := store.HGetVal(KeyAccount, accountID, account); ok {
		return account, nil
	}
	if account, err = accountModel.GetAccount(accountID); err != nil {
		return nil, err
	}
	if err := store.HSetVal(KeyAccount, account.ID, account); err != nil {
		log.Error("HSet failed:", err)
	}
	return account, err
}
