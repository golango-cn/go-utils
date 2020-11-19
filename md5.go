package go_utils

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
)

// Md5加密
func Md5(str string) string {
	//方法二
	w := md5.New()
	_, _ = io.WriteString(w, str)
	//将str写入到w中
	v := fmt.Sprintf("%x", w.Sum(nil))
	return v
}

// Sha1加密
func Sha1(str string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(str))
	v := fmt.Sprintf("%x", h.Sum(nil))
	return v
}
