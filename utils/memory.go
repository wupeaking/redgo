package utils

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Malloc 分配内存
func Malloc(size int64) unsafe.Pointer {
	return unsafe.Pointer(C.malloc(C.size_t(size)))
}

// Free 释放内存
func Free(ptr unsafe.Pointer) {
	C.free(ptr)
}
