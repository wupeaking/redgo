package datastruct

import (
	"errors"
	"unsafe"
)

// ListProcess 定义list的函数族
type ListProcess interface {
	ListDup(*ListNode) (*ListNode, error) // 复制一个链表节点所保存的值
	ListFree(*ListNode) error             // 释放链表节点所保存的值
	ListMatch(l1, l2 *ListNode) bool      // 用于对比链表节点所保存的值是否和输入值相等
}

// List 声明列表对象
type List struct {
	head      *ListNode // 指向列表的头部
	tail      *ListNode // 指向列表的尾部
	len       int64     // 列表的长度
	Processor ListProcess
}

// ListNode 列表节点
type ListNode struct {
	Prev  *ListNode
	Next  *ListNode
	Value interface{}
}

// ListCreate 创建一个不包含任何节点的新链表
func ListCreate() *List {
	return new(List)
}

// ListSetProcess 设置链表的操作函数
func (l *List) ListSetProcess(proc ListProcess) {
	l.Processor = proc
}

// ListLengeth 返回列表的长度
func (l *List) ListLengeth() int64 {
	return l.len
}

// ListFirst 返回列表的表头
func (l *List) ListFirst() *ListNode {
	return l.head
}

// ListLast 返回列表的尾部
func (l *List) ListLast() *ListNode {
	return l.tail
}

// ListPrevNode 返回给定节点的前一个节点
func (l *List) ListPrevNode(node *ListNode) *ListNode {
	if node == nil {
		return nil
	}
	return node.Prev
}

// ListNextNode 返回给定节点的下一个节点
func (l *List) ListNextNode(node *ListNode) *ListNode {
	if node == nil {
		return nil
	}
	return node.Next
}

// ListNodeValue 返回指定节点保存的值
func (l *List) ListNodeValue(node *ListNode) interface{} {
	if node == nil {
		return nil
	}
	return node.Value
}

// ListAddNodeHead 将一个指定的节点加入表头
func (l *List) ListAddNodeHead(node *ListNode) bool {
	if node == nil {
		return false
	}
	headNode := l.head
	if headNode == nil {
		// 说明这个链表还没有节点
		l.head = node
		node.Prev = nil
		node.Next = nil
		l.tail = node
		l.len++
		return true
	}
	headNode.Prev = node
	node.Next = headNode
	node.Prev = nil
	l.head = node
	l.len++
	return true
}

// ListAddNodeTail 将一个节点接入表尾
func (l *List) ListAddNodeTail(node *ListNode) bool {
	if node == nil {
		return false
	}
	tailNode := l.tail
	if tailNode == nil {
		// 说明这个链表还没有节点
		l.tail = node
		node.Prev = nil
		node.Next = nil
		l.head = node
		l.len++
		return true
	}
	tailNode.Next = node
	node.Prev = tailNode
	node.Next = nil
	l.tail = node
	l.len++
	return true
}

// ListInsertNode 将一个新节点插入指定节点的前或者尾 dir: 负数或0表示向前插入 正数表示向后插入
func (l *List) ListInsertNode(newNode *ListNode, target *ListNode, dir int) (bool, error) {
	if newNode == nil {
		return false, errors.New("新节点地址为空")
	}
	t := unsafe.Pointer(target)
	// 从头开始遍历
	ptr := l.head
	for ptr != nil {
		if t != unsafe.Pointer(ptr) {
			ptr = ptr.Next
			continue
		}
		// 执行到此处说明找到了定位节点
		if dir <= 0 {
			// 在这个节点前方加入节点 需要判断这个节点是不是首节点
			prev := ptr.Prev
			if prev == nil {
				// 说明此节点是一个首节点
				ptr.Prev = newNode
				newNode.Next = ptr
				newNode.Prev = nil
				l.head = newNode
			} else {
				prev.Next = newNode
				newNode.Prev = prev
				newNode.Next = ptr
				ptr.Prev = newNode
			}
		} else {
			next := ptr.Next
			if next == nil {
				ptr.Next = newNode
				newNode.Prev = ptr
				newNode.Next = nil
				l.tail = newNode
			} else {
				ptr.Next = newNode
				newNode.Prev = ptr
				newNode.Next = next
				next.Prev = newNode
			}
		}
		l.len++
		return true, nil
	}
	return false, errors.New("未找到定位节点")
}

// ListInsertNodeByValue 将一个新节点插入指定节点的前或者尾 (比较节点的值是否相等 而不是地址) dir: 负数或0表示向前插入 正数表示向后插入
func (l *List) ListInsertNodeByValue(newNode *ListNode, target *ListNode, dir int) (int, error) {
	if newNode == nil {
		return 0, errors.New("新节点地址为空")
	}
	// 从头开始遍历
	i := 0
	ptr := l.head
	for ptr != nil {
		if !l.Processor.ListMatch(target, ptr) {
			i++
			ptr = ptr.Next
			continue
		}
		i++
		// 执行到此处说明找到了定位节点
		if dir <= 0 {
			// 在这个节点前方加入节点 需要判断这个节点是不是首节点
			prev := ptr.Prev
			if prev == nil {
				// 说明此节点是一个首节点
				ptr.Prev = newNode
				newNode.Next = ptr
				newNode.Prev = nil
				l.head = newNode
			} else {
				prev.Next = newNode
				newNode.Prev = prev
				newNode.Next = ptr
				ptr.Prev = newNode
			}
		} else {
			next := ptr.Next
			if next == nil {
				ptr.Next = newNode
				newNode.Prev = ptr
				newNode.Next = nil
				l.tail = newNode
			} else {
				ptr.Next = newNode
				newNode.Prev = ptr
				newNode.Next = next
				next.Prev = newNode
			}
		}
		l.len++
		return i, nil
	}
	return i, errors.New("为找到定位节点")
}

// ListInsterNodeByIndex 通过索引插入一个节点
func (l *List) ListInsterNodeByIndex(newNode *ListNode, index int64, dir int) (bool, error) {
	if index < 0 || newNode == nil || index >= l.len {
		return false, errors.New("传递参数错误")
	}
	ptr := l.head
	i := int64(0)
	// todo:: 此处需要优化 当索引值过半时  从后往前检索
	for i != index {
		ptr = ptr.Next
		i++
	}
	if dir <= 0 {
		// 在这个节点前方加入节点 需要判断这个节点是不是首节点
		prev := ptr.Prev
		if prev == nil {
			// 说明此节点是一个首节点
			ptr.Prev = newNode
			newNode.Next = ptr
			newNode.Prev = nil
			l.head = newNode
		} else {
			prev.Next = newNode
			newNode.Prev = prev
			newNode.Next = ptr
			ptr.Prev = newNode
		}
	} else {
		next := ptr.Next
		if next == nil {
			ptr.Next = newNode
			newNode.Prev = ptr
			newNode.Next = nil
			l.tail = newNode
		} else {
			ptr.Next = newNode
			newNode.Prev = ptr
			newNode.Next = next
			next.Prev = newNode
		}
	}
	l.len++
	return true, nil
}

// ListIndex 查找指定索引的节点
func (l *List) ListIndex(index int64) *ListNode {
	if index < 0 || index >= l.len {
		return nil
	}
	ptr := l.head
	i := int64(0)
	// todo:: 此处需要优化 当索引值过半时  从后往前检索
	for i != index {
		i++
		ptr = ptr.Next
	}
	return ptr
}

// ListDelNodeByIndex 根据索引删除指定节点
func (l *List) ListDelNodeByIndex(index int64) {
	if index < 0 || index >= l.len {
		return
	}
	ptr := l.head
	i := int64(0)
	// todo:: 此处需要优化 当索引值过半时  从后往前检索
	for i != index {
		ptr = ptr.Next
		i++
	}
	// 判断是否是首节点
	prev := ptr.Prev
	next := ptr.Next
	if prev == nil && next != nil {
		// 只是首节点
		l.head = next
		next.Prev = nil
	} else if prev != nil && next == nil {
		// 只是未节点
		prev.Next = nil
		l.tail = prev
	} else if prev == nil && next == nil {
		// 既是首节点又是尾节点
		l.head = nil
		l.tail = nil
	} else {
		// 既不是首节点也不是未节点
		prev.Next = next
		next.Prev = prev
	}
	l.len--
}

// ListDelNodeByValue 根据值相等来删除指定节点
func (l *List) ListDelNodeByValue(valueNode *ListNode, dir int) bool {
	// 从头开始遍历
	var ptr *ListNode
	if dir > 0 {
		ptr = l.head
	} else {
		ptr = l.tail
	}

	for ptr != nil {
		if !l.Processor.ListMatch(valueNode, ptr) {
			if dir > 0 {
				ptr = ptr.Next
			} else {
				ptr = ptr.Prev
			}
			continue
		}
		// 执行到此处说明找到了定位节点
		// 判断是否是首节点
		prev := ptr.Prev
		next := ptr.Next
		if prev == nil && next != nil {
			// 只是首节点
			l.head = next
			next.Prev = nil
		} else if prev != nil && next == nil {
			// 只是未节点
			prev.Next = nil
			l.tail = prev
		} else if prev == nil && next == nil {
			// 既是首节点又是尾节点
			l.head = nil
			l.tail = nil
		} else {
			// 既不是首节点也不是未节点
			prev.Next = next
			next.Prev = prev
		}
		l.len--
		return true
	}
	return false
}

// ListTrim 剔除范围外的列表
func (l *List) ListTrim(left int, right int) {
	ptr := l.head
	cur := 0
	isStart := false
	isEnd := false
	for ptr != nil {
		tmp := ptr.Next
		if cur == left && isStart != true {
			l.head = ptr
			ptr.Prev = nil
		}
		if cur == right && isEnd != true {
			l.tail = ptr
			ptr.Next = nil
		}
		cur++
		if isStart && isEnd {
			break
		}
		ptr = tmp
	}
	l.len = int64(right - left + 1)
}

// ListCopy 返回一个链表额副本
func (l *List) ListCopy() (*List, error) {
	newlist := ListCreate()
	ptr := l.head
	for ptr != nil {
		node, e := l.Processor.ListDup(ptr)
		if e != nil {
			return nil, e
		}
		newlist.ListAddNodeTail(node)
		ptr = ptr.Next
	}
	newlist.len = l.len
	return newlist, nil
}

// ListFree 释放链表空间
// todo:: 此处有疑问 如果我只取消节点的引用 gc释放能够自动释放内存 现在先手动释放 下回有时间检测一下
func (l *List) ListFree() {
	ptr := l.head
	for ptr != nil {
		next := ptr.Next
		ptr.Prev = nil
		ptr.Next = nil
		ptr.Value = nil
		if next == nil {
			break
		}
		ptr = next
	}
	l.head = nil
	l.tail = nil
	l.len = 0
}

// StrListNode 字符串类型的函数簇
type StrListNode struct{}

//StrListFuncs 字符串类型的列表函数簇
var StrListFuncs = &StrListNode{}

// ListDup 拷贝
func (*StrListNode) ListDup(node *ListNode) (*ListNode, error) {
	value := node.Value.(string)
	result := new(ListNode)
	result.Value = value
	return result, nil
}

// ListFree 释放
func (*StrListNode) ListFree(node *ListNode) error {
	// prev := node.Prev
	// next := node.Next
	// prev.Next = next
	// next.Prev = prev
	node.Value = nil
	return nil
}

// ListMatch 比较
func (*StrListNode) ListMatch(n1, n2 *ListNode) bool {
	return n1.Value.(string) == n2.Value.(string)
}
