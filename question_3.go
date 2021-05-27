package stream_test

import (
	"strings"

	"stream_test/stream"
)

// Question3Sub1 Q1: 输入一个整数 int，字符串string。将这个字符串重复n遍返回
func Question3Sub1(str string, n int) string {
	return stream.Repeat(str, n).Reduce(func(acc stream.R, e stream.T, idx int, sLen int) stream.R {
		acc.(*strings.Builder).WriteString(e.(string))
		return acc
	}, &strings.Builder{}).(*strings.Builder).String()
}
