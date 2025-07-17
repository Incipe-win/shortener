package base62

import (
	"math"
	"slices"
	"strings"
)

// 62 进制转换的模块

var (
	base62Str string
)

// MustInit 要使用base62这包必须要调用该函数完成初始化
func MustInit(bs string) {
	if len(bs) == 0 {
		panic("base62Str is empty")
	}
	base62Str = bs
}

// Int2String 10 进制转成 62 进制
func Int2String(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	bl := []byte{}
	for seq > 0 {
		mod := seq % 62
		div := seq / 62
		bl = append(bl, base62Str[mod])
		seq = div
	}
	slices.Reverse(bl)
	return string(bl)
}

// String2Int 62进制字符串转为10进制数
func String2Int(s string) (seq uint64) {
	bl := []byte(s)
	slices.Reverse(bl)
	for idx, b := range bl {
		base := math.Pow(62, float64(idx))
		seq += uint64(strings.Index(base62Str, string(b))) * uint64(base)
	}
	return seq
}
