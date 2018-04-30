package db

import (
	"errors"
)

var (
	ErrKeyDuplicate = errors.New("key duplicate")
)
