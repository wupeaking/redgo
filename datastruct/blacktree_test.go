package datastruct

import "testing"
import "strconv"

func TestBlackTree(t *testing.T) {
	// 新建红黑树对象
	rbt := NewRbt()
	// 创建节点
	for i := 0; i < 10; i++ {
		node := NewRBTNode(StringNodeFuncs, i, strconv.Itoa(i))
		rbt.Insert(node)
	}
	for i := 0; i < 10; i++ {
		node, ok := rbt.Get(i)
		if !ok {
			t.Fatal("获取节点失败")
		}
		score, v := node.GetNodeInfo()
		if score != i || v.(string) != strconv.Itoa(i) {
			t.Fatal("获取节点内容失败")
		}
	}
	// 测试获取的最大值节点
	node, ok := rbt.GetMaxNode()
	if !ok {
		t.Fatal("获取最大值节点失败")
	}
	score, v := node.GetNodeInfo()
	if score != 9 || v.(string) != "9" {
		t.Fatal("获取最大值节点错误")
	}

	// 测试获取最小值节点
	node, ok = rbt.GetMaxNode()
	if !ok {
		t.Fatal("获取最大值节点失败")
	}
	score, v = node.GetNodeInfo()
	if score != 9 || v.(string) != "9" {
		t.Fatal("获取最大值节点错误")
	}

	//
	_, ok = rbt.Get(10)
	if ok {
		t.Fatal("获取到错误节点")
	}
}
