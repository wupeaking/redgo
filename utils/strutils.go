package utils

import "unsafe"

// Str2bytes 字符串转为字节
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//Bytes2str 字节转为字符串
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
