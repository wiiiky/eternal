package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	mrand "math/rand"
	"time"
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

func RandDigit(size int) string {
	rnd := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}
