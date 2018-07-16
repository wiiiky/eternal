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
	expire := time.Duration(0)
	if !token.ETime.IsZero() {
		if expire = token.ETime.Sub(time.Now()); expire == 0 {
			/* 临界情况，如果token有有效期，且有效期正好等于当前时间，此时不设置缓存 */
			expire = -1
		}
	}
	if expire >= 0 { /* 只有当token永久有效，或者处于有效期中时才设置缓存 */
		if err = store.SetVal(key, token, expire); err != nil {
			log.Error("SetVal failed:", err)
		}
	}
	return token, nil
}

func DeleteToken(tokenID string) error {
	if err := accountModel.DeleteToken(tokenID); err != nil {
		return err
	}
	key := getTokenKey(tokenID)
	store.Del(key)
	return nil
}
