package account

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func encryptPasswordWithSHA1(salt, plain string) string {
	var buf bytes.Buffer
	buf.WriteByte('^')
	buf.WriteString(salt)
	buf.WriteByte(':')
	buf.WriteString(plain)
	buf.WriteByte('$')

	data := sha1.Sum(buf.Bytes())
	return hex.EncodeToString(data[:])
}

func encryptPasswordWithMD5(salt, plain string) string {
	var buf bytes.Buffer
	buf.WriteByte('^')
	buf.WriteString(plain)
	buf.WriteByte(':')
	buf.WriteString(salt)
	buf.WriteByte('$')

	data := md5.Sum(buf.Bytes())
	return hex.EncodeToString(data[:])
}

func encryptPasswordWithSHA256(salt, plain string) string {
	var buf bytes.Buffer
	buf.WriteByte('$')
	buf.WriteString(salt)
	buf.WriteByte('.')
	buf.WriteString(plain)
	buf.WriteByte('^')

	data := md5.Sum(buf.Bytes())
	return hex.EncodeToString(data[:])
}

/*
 * 对明文密码进行加密
 * 参数加密salt、明文密码、加密方式
 */
func encryptPassword(salt, plain, ptype string) (string, string) {
	switch ptype {
	case PTYPE_SHA1:
		return encryptPasswordWithSHA1(salt, plain), ptype
	case PTYPE_MD5:
		return encryptPasswordWithMD5(salt, plain), ptype
	default:
		return encryptPasswordWithSHA256(salt, plain), PTYPE_SHA256
	}
}
