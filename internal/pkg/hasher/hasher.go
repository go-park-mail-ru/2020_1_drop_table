package hasher

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetSha1(value []byte) string {
	sum := sha1.Sum(value)
	return hex.EncodeToString(sum[:])
}
