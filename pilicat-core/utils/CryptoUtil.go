package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

//生成32位md5字串
func GetMd5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetSaltMD5(input, salt string) string {
	hash := md5.New()
	//salt = "salt123456" //盐值
	io.WriteString(hash, input+salt)
	result := fmt.Sprintf("%x", hash.Sum(nil))
	return result
}

//生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5(base64.URLEncoding.EncodeToString(b))
}
