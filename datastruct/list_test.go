package datastruct

import "testing"
import "strconv"

type Str struct{}

func (*Str) ListDup(node *ListNode) (*ListNode, error) {
	value := node.Value.(string)
	result := new(ListNode)
	result.Value = value
	return result, nil
}

func (*Str) ListFree(node *ListNode) error {
	// prev := node.Prev
	// next := node.Next
	// prev.Next = next
	// next.Prev = prev
	node.Value = nil
	return nil
}

func (*Str) ListMatch(n1, n2 *ListNode) bool {
	return n1.Value.(string) == n2.Value.(string)
}

func TestList(t *testing.T) {
	l := ListCreate()
	l.ListSetProcess(&Str{})
	// 创建几个listnode
	for i := 0; i < 10; i++ {
		node := ListNode{Value: strconv.Itoa(i)}
		l.ListAddNodeTail(&node)
	}
	if l.ListLengeth() != 10 {
		t.Fatal("创建索引错误", l.ListLengeth())
	}
	// 索引迭代
	for i := 0; i < 10; i++ {
		if l.ListIndex(int64(i)).Value.(string) != strconv.Itoa(i) {
			t.Fatal("索引迭代链表出现异常")
		}
	}
	// 索引删除
	for i := 0; i < 10; i++ {
		l.ListDelNodeByIndex(int64(0))
	}
	if l.ListLengeth() != 0 || l.ListFirst() != nil || l.ListLast() != nil {
		t.Fatal("索引删除出现异常", l.ListLengeth())
	}

	// 重新创建
	for i := 0; i < 10; i++ {
		node := ListNode{Value: strconv.Itoa(i)}
		l.ListAddNodeTail(&node)
	}
	//在索引头部添加一个节点
	l.ListAddNodeHead(&ListNode{Value: "-1"})
	if l.len != 11 || l.ListFirst().Value.(string) != "-1" {
		t.Fatal("头部添加索引出现异常")
	}
	// 在表未添加
	l.ListAddNodeTail(&ListNode{Value: "10"})
	if l.len != 12 || l.ListLast().Value.(string) != "10" {
		t.Fatal("头部添加索引出现异常")
	}

	// 找到11这个节点
	node := l.ListIndex(int64(11))
	if node.Value.(string) != "10" {
		t.Fatal("索引指定节点出现错误")
	}
	l.ListInsertNode(&ListNode{Value: "before"}, node, -1)
	if l.ListIndex(int64(11)).Value.(string) != "before" {
		t.Fatal("指定节点前向插入错误")
	}
	l.ListInsertNode(&ListNode{Value: "after"}, node, 1)
	if l.ListIndex(int64(13)).Value.(string) != "after" {
		t.Fatal("指定节点向后插入错误", l.ListIndex(int64(13)).Value.(string))
	}

	// 删除掉刚才创建的节点
	l.ListDelNodeByIndex(int64(13))
	l.ListDelNodeByIndex(int64(11))
	if l.ListLengeth() != 12 {
		t.Fatal("删除节点异常")
	}

	// 拷贝一个新链表
	newList, e := l.ListCopy()
	if e != nil || newList.ListLengeth() != 12 {
		t.Fatal("拷贝新链表失败", e)
	}
	// 迭代出来
	for i := int64(0); i < newList.ListLengeth(); i++ {
		if newList.ListIndex(i).Value.(string) != strconv.Itoa(int(i-1)) {
			t.Fatal("链表拷贝出现错误")
		}
		t.Log("index: ", i, "value: ", newList.ListIndex(i).Value.(string))
	}

}
