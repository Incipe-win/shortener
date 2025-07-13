package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// Sum 对传入的参数求 md5 值
func Sum(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
