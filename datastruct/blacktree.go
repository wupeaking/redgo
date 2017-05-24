package datastruct

// 红黑树的实现 用于实现排序操作
import (
	"errors"

	rbt "github.com/wupeaking/gods/trees/redblacktree"
)

// RBT 红黑树对象
type RBT struct {
	tree *rbt.Tree
}

// RBTNode 每个红黑树的节点
type RBTNode struct {
	score int         // 分数 用于排序
	value interface{} // 用于保存值
	funcs RBTNodeFuncs
}

// RBTNodeFuncs 聚合函数
type RBTNodeFuncs interface {
	// 更改节点的值
	SetValue(node *RBTNode, value interface{}) error
	// 删除node中为value的值
	DeleteValue(node *RBTNode, value interface{}) error
}

// NewRbt 新建一个红黑树操作对象
func NewRbt() *RBT {
	o := new(RBT)
	o.tree = rbt.NewWithIntComparator()
	return o
}

// NewRBTNode 新的红黑树节点
func NewRBTNode(funcs RBTNodeFuncs, socre int, value interface{}) *RBTNode {

	node := &RBTNode{score: socre, funcs: funcs}
	node.funcs.SetValue(node, value)
	return node
}

// --------节点相关操作--------------

// SetNodeValue 设置node的值
func (node *RBTNode) SetNodeValue(value interface{}) error {
	return node.funcs.SetValue(node, value)
}

// DeleteNodeValue 删除node中的一个值
func (node *RBTNode) DeleteNodeValue(value interface{}) error {
	return node.funcs.DeleteValue(node, value)
}

// GetNodeInfo 返回节点信息 分数 值
func (node *RBTNode) GetNodeInfo() (score int, value interface{}) {
	return node.score, node.value
}

// -------树相关操作

// Insert 将节点加入红黑树中
func (t *RBT) Insert(node *RBTNode) {
	t.tree.Put(node.score, node)
}

// Delete 将节点从树中删除
func (t *RBT) Delete(node *RBTNode) {
	t.tree.Remove(node.score)
}

// Get 获取一个节点的内容
func (t *RBT) Get(score int) (*RBTNode, bool) {
	value, found := t.tree.Get(score)
	if !found {
		return nil, false
	}
	return value.(*RBTNode), true
}

// GetMinNode 获取最小节点
func (t *RBT) GetMinNode() (*RBTNode, bool) {
	value := t.tree.Left().Value
	if value == nil {
		return nil, false
	}
	return value.(*RBTNode), true
}

// GetMaxNode 获取最小节点
func (t *RBT) GetMaxNode() (*RBTNode, bool) {
	value := t.tree.Right().Value
	if value == nil {
		return nil, false
	}
	return value.(*RBTNode), true
}

//StringRBTNode 字符串类型节点对象操作函数簇
type StringRBTNode struct {
}

// StringNodeFuncs 字符串类型节点的操作函数
var StringNodeFuncs = &StringRBTNode{}

//SetValue 设置节点值
func (*StringRBTNode) SetValue(node *RBTNode, value interface{}) error {
	if _, ok := value.(string); !ok {
		return errors.New("value must be string")
	}
	node.value = value
	return nil
}

// DeleteValue 删除一个值
func (*StringRBTNode) DeleteValue(node *RBTNode, value interface{}) error {
	node.value = nil
	return nil
}
