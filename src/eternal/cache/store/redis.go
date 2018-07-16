package store

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

var _redis *redis.Client

func REDIS() *redis.Client {
	return _redis
}

func InitRedis(sURL string) error {
	opt, err := redis.ParseURL(sURL)
	if err != nil {
		return err
	}
	_redis = redis.NewClient(opt)
	_, err = _redis.Ping().Result()
	return err
}

func SetVal(k string, v interface{}, etime time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return _redis.Set(k, data, etime).Err()
}

func HSetVal(k, f string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return _redis.HSet(k, f, data).Err()
}

func GetVal(k string, v interface{}) (bool, error) {
	val, err := _redis.Get(k).Result()
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(val), v); err != nil {
		return false, nil
	}
	return true, nil
}

func HGetVal(k, f string, v interface{}) (bool, error) {
	val, err := _redis.HGet(k, f).Result()
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(val), v); err != nil {
		return false, nil
	}
	return true, nil
}

func Del(keys ...string) error {
	return _redis.Del(keys...).Err()
}
