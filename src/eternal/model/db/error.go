package db

import (
	"errors"
)

var (
	ErrKeyDuplicate = errors.New("key duplicate")
	ErrKeyNotFound  = errors.New("key not found")
)
