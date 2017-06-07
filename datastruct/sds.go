package datastruct

// 构造redis的二进制安全的动态字符串结构(SDS)

var lenThreshold = 1024 * 1024

// Sds -- sds 结构体
type Sds struct {
	len  int    // 记录当前使用的字符串长度
	free int    // 记录还剩下的个数
	buf  []byte // 保存字符内容
}

// NewSds -- 创建一个包含字符串的新的sds对象
func NewSds(str []byte) *Sds {
	strlen := len(str)
	if strlen <= lenThreshold {
		newbuf := make([]byte, 2*strlen, 2*strlen)
		copy(newbuf, str)
		return &Sds{len: strlen, free: strlen, buf: newbuf}
	}
	newbuf := make([]byte, strlen+lenThreshold, strlen+lenThreshold)
	copy(newbuf, str)
	return &Sds{len: strlen, free: lenThreshold, buf: newbuf}
}

//SdsEmpty -- 返回一个空的sds对象
func SdsEmpty() *Sds {
	return &Sds{len: 0, free: 0, buf: make([]byte, 0, 0)}
}

// SdsFree -- 释放一个sds的内存
func (sds *Sds) SdsFree() {
	sds.buf = nil
	sds.len = 0
	sds.free = 0
}

// SdsLen -- 返回字符串的长度
func (sds *Sds) SdsLen() int {
	return sds.len
}

// SdsVail --返回已经使用的长度
func (sds *Sds) SdsVail() int {
	return sds.free
}

// SdsCopy 创建一个新的副本
func (sds *Sds) SdsCopy() *Sds {
	c := new(Sds)
	c.len = sds.len
	c.free = sds.free
	c.buf = make([]byte, c.len+c.free, c.len+c.free)
	copy(c.buf, sds.buf)
	return c
}

// SdsClear 清空内容
func (sds *Sds) SdsClear() {
	// 如果长度大于设定的阈值 则释放
	if sds.len > lenThreshold {
		sds.len = 0
		sds.free = 0
		sds.buf = nil
		return
	}
	sds.free = sds.len
	sds.len = 0
	return
}

// SdsCat 将指定的字符串拼接到末尾
func (sds *Sds) SdsCat(str []byte) {
	strlen := len(str)
	if sds.free > strlen {
		copy(sds.buf[sds.len:], str)
		sds.free -= strlen
		sds.len += strlen
		return
	}
	// 如果剩余的空间不够 需要进行扩容 根据条件进行扩容
	size := 0
	if sds.len+strlen > lenThreshold {
		// 如果长度大于阈值 只多分配1M
		size = sds.len + strlen + lenThreshold
	} else {
		size = 2 * (sds.len + strlen)
	}
	newbuf := make([]byte, size, size)
	copy(newbuf, sds.buf[0:sds.len])
	copy(newbuf[sds.len:], str)
	sds.len = sds.len + strlen
	sds.free = size - sds.len
	sds.buf = newbuf
	return
}

// SdsCatSds 将给定的sds拼接到一个sds的末尾
func (sds *Sds) SdsCatSds(newSds *Sds) {
	size := 0
	if sds.len+newSds.len > lenThreshold {
		size = sds.len + newSds.len + lenThreshold
	} else {
		size = 2 * (sds.len + newSds.len)
	}
	newbuf := make([]byte, size, size)
	copy(newbuf, sds.buf[0:sds.len])
	copy(newbuf[sds.len:], newSds.buf[0:newSds.len])
	sds.len = sds.len + newSds.len
	sds.free = size - sds.len
	sds.buf = newbuf
	return
}

// SdsCopyStr 拷贝一个字符串 覆盖原来的内容
func (sds *Sds) SdsCopyStr(str []byte) {
	strlen := len(str)
	if sds.len+sds.free >= strlen {
		copy(sds.buf, str)
		sds.free = sds.len + sds.free - strlen
		sds.len = strlen
		return
	}
	size := 0
	if strlen > lenThreshold {
		size = strlen + lenThreshold
	} else {
		size = 2 * strlen
	}
	newbuf := make([]byte, size, size)
	copy(newbuf, str)
	sds.len = strlen
	sds.free = size - sds.len
	sds.buf = newbuf
}

// String 以字符串格式返回内容
func (sds *Sds) String() string {
	return string(sds.buf[0:sds.len])
}

// Buffer 以字节返回
func (sds *Sds) Buffer() []byte {
	return sds.buf[0:sds.len]
}

// SdsRange 保留指定返回的内容
func (sds *Sds) SdsRange(left, right int) *Sds {
	if right > sds.len {
		right = sds.len
	}
	if left > sds.len {
		left = sds.len
	}
	size := right - left
	newbuf := make([]byte, size, size)
	copy(newbuf, sds.buf[left:right])
	sds.len = size
	sds.free = 0
	sds.buf = newbuf
	return sds
}

// SdsCmp 比较两个sds的字符串是否一样
func (sds *Sds) SdsCmp(other *Sds) bool {
	if sds.len != sds.len {
		return false
	}
	compare := func(a, b []byte) bool {
		for i, v := range a {
			if v != b[i] {
				return false
			}
			continue
		}
		return true
	}
	return compare(sds.buf[0:sds.len], other.buf[0:other.len])
}
