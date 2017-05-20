package datastruct

// 这是一个测试文件 测试sds的功能是否正确

import (
	"testing"
)

func TestNewSds(t *testing.T) {
	sds := NewSds([]byte("this is test"))
	if sds.len != len("this is test") || sds.free != len("this is test") || sds.String() != "this is test" {
		t.Fatal(" 创建sds对象出现问题", sds.SdsLen(), sds.SdsVail(), sds.String())
	} else {
		t.Log("创建sds对象测试成功")
	}

	// 拼接字符串
	sds.SdsCat([]byte(" new str"))
	if sds.String() != "this is test new str" {
		t.Fatal("拼接字符串失败")
	}
	if sds.SdsLen() != len("this is test new str") {
		t.Fatal("返回使用空间错误")
	}
	if sds.SdsVail() != len("this is test")-len(" new str") {
		t.Fatal("返回剩余空间错误", sds.SdsVail())
	}

	// 拼接 sds
	newsds := NewSds([]byte(" hello world"))
	sds.SdsCatSds(newsds)
	if sds.String() != "this is test new str hello world" {
		t.Fatal("拼接字符串失败")
	}
	if sds.SdsLen() != len("this is test new str hello world") {
		t.Fatal("返回使用空间错误")
	}
	if sds.SdsVail() != len("this is test new str hello world") {
		t.Fatal("返回剩余空间错误")
	}

	// 比较字符串
	if !sds.SdsCmp(NewSds([]byte("this is test new str hello world"))) {
		t.Fatal("字符串比较错误")
	}
	// 获取区间字符
	if sds.SdsRange(5, 12).String() != "is test" {
		t.Fatal("获取区间字符串错误")
	}
	// 创建一个副本
	copySds := sds.SdsCopy()
	if !sds.SdsCmp(sds) || sds.SdsVail() != copySds.SdsVail() {
		t.Fatal("创建副本失败")
	}
	// 清空字符串
	copySds.SdsClear()
	if copySds.SdsLen() != 0 || copySds.SdsVail() != sds.SdsLen() {
		t.Fatal("清空字符串错误")
	}
	// 是否内容
	copySds.SdsFree()
	if copySds.SdsLen() != 0 || copySds.SdsVail() != 0 || copySds.buf != nil {
		t.Fatal("释放sds出现错误")
	}
}

func TestEmptySds(t *testing.T) {
	sds := SdsEmpty()
	if sds.free != 0 || sds.len != 0 {
		t.Fatal("创建空的sds失败")
	}
	// 拷贝字符串
	sds.SdsCopyStr([]byte("new str"))
	if sds.SdsLen() != len("new str") || sds.SdsVail() != len("new str") || sds.String() != "new str" {
		t.Fatal("拷贝字符串出现错误")
	}
}
