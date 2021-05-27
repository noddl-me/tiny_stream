package stream_test

import (
	"unsafe"

	"stream_test/stream"
)

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Question2Sub1 Q1: 计算一个 string 中小写字母的个数
func Question2Sub1(str string) int64 {
	return int64(stream.OfSlice(str).
		Filter(func(e stream.T) bool {
			return e.(byte) >= 97 && e.(byte) <= 123
		}).Count())
}

// Question2Sub2 Q2: 找出 []string 中，包含小写字母最多的字符串
func Question2Sub2(list []string) string {
	return stream.OfSlice(list).
		Map(func(e stream.T) stream.R {
			return &stream.Pair{
				First:  e,
				Second: Question2Sub1(e.(string)),
			}
		}).
		Reduce(func(acc stream.R, e stream.T, idx int, sLen int) stream.R {
			if e.(*stream.Pair).Second.(int64) > acc.(*stream.Pair).Second.(int64) {
				return e
			}
			return acc
		}, &stream.Pair{First: "", Second: int64(0)}).(*stream.Pair).First.(string)
}
