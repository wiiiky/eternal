package util

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func RandString(size int) (string, error) {
	bl := size / 2
	if size%2 != 0 {
		bl++
	}

	buf := make([]byte, bl)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
