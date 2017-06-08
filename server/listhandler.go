package server

// list相关命令的处理
import (
	"errors"

	"github.com/wupeaking/redgo/datastruct"
)

// LPush 在列表的头部添加数据
func (myself *SrvHandler) LPush(key string, values [][]byte) (int, error) {
	if len(values) == 0 {
		return 0, errors.New("params erros")
	}

	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		list := datastruct.ListCreate()
		list.ListSetProcess(datastruct.StrListFuncs)
		for _, value := range values {
			listNode := &datastruct.ListNode{Value: string(value)}
			list.ListAddNodeHead(listNode)
		}

		v := &Value{value: list, valueType: LIST}
		myself.db.data.Set(key, v)
		return int(list.ListLengeth()), nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != LIST {
		return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	listValue := comValue.value.(*datastruct.List)
	for _, value := range values {
		listNode := &datastruct.ListNode{Value: string(value)}
		listValue.ListAddNodeHead(listNode)
	}

	return int(listValue.ListLengeth()), nil
}

// RPush 在尾部加入数据
func (myself *SrvHandler) RPush(key string, values [][]byte) (int, error) {
	if len(values) == 0 {
		return 0, errors.New("params erros")
	}
	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		list := datastruct.ListCreate()
		list.ListSetProcess(datastruct.StrListFuncs)
		for _, value := range values {
			listNode := &datastruct.ListNode{Value: string(value)}
			list.ListAddNodeTail(listNode)
		}

		v := &Value{value: list, valueType: LIST}
		myself.db.data.Set(key, v)
		return int(list.ListLengeth()), nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != LIST {
		return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	listValue := comValue.value.(*datastruct.List)
	for _, value := range values {
		listNode := &datastruct.ListNode{Value: string(value)}
		listValue.ListAddNodeTail(listNode)
	}

	return int(listValue.ListLengeth()), nil
}

// LInsert 在某个值的前后插入新值
func (myself *SrvHandler) LInsert(key string, dir string, pivot string, value []byte) (int, error) {
	//todo :: 需要实现
	return 0, nil
}

// LSet 修改指定下标的值
func (myself *SrvHandler) LSet(key string, index int, value []byte) error {
	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		return errors.New("not exist the key")
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != LIST {
		return errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	listValue := comValue.value.(*datastruct.List)
	size := int(listValue.ListLengeth())
	//判断 参数是否符合要求
	if index < 0 {
		index = index + size
	}
	// 如果此时index 小于0 或者 大于等于最大长度 则为超出范围
	if index < 0 || index >= size {
		return errors.New("ERR index out of range")
	}
	listNode := &datastruct.ListNode{Value: string(value)}
	_, e := listValue.ListInsterNodeByIndex(listNode, int64(index), 1)
	return e
}

// 读取相关----------

// LRange 读取list某个范围的元素
func (myself *SrvHandler) LRange(key string, left int, right int) ([][]byte, error) {
	//判断 参数是否符合要求
	if left < 0 {
		return nil, errors.New("left must be a positive numberand")
	}

	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		return nil, nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != LIST {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	listValue := comValue.value.(*datastruct.List)
	size := int(listValue.ListLengeth())
	absright := 0
	// 确定右坐标
	if right < 0 {
		absright = size + right + 1
	} else {
		absright = right
	}
	startNode := listValue.ListIndex(int64(left))
	result := make([][]byte, 0)

	for i := 0; i < absright-left; i++ {
		if startNode == nil {
			break
		}
		result = append(result, []byte(startNode.Value.(string)))
		startNode = startNode.Next
	}
	return result, nil
}

// LLen 返回列表长度
func (myself *SrvHandler) LLen(key string) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		return 0, nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != LIST {
		return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	listValue := comValue.value.(*datastruct.List)

	return int(listValue.ListLengeth()), nil
}

// LIndex 获取指定下标的元素
func (myself *SrvHandler) LIndex(key string, index int) ([]byte, error) {
	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		return nil, nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != LIST {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	listValue := comValue.value.(*datastruct.List)
	node := listValue.ListIndex(int64(index))
	if node == nil {
		return nil, nil
	}
	return []byte(node.Value.(string)), nil
}

// 删除相关-------

// LPop 从左边删除一个元素
func (myself *SrvHandler) LPop(key string) ([]byte, error) {
	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		return nil, nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != LIST {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	listValue := comValue.value.(*datastruct.List)
	node := listValue.ListFirst()
	if node == nil {
		return nil, nil
	}
	listValue.ListDelNodeByIndex(int64(0))
	return []byte(node.Value.(string)), nil
}

// RPop 从尾部删除一个元素
func (myself *SrvHandler) RPop(key string) ([]byte, error) {
	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		return nil, nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != LIST {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	listValue := comValue.value.(*datastruct.List)
	node := listValue.ListLast()
	if node == nil {
		return nil, nil
	}
	size := listValue.ListLengeth()
	listValue.ListDelNodeByIndex(int64(size - 1))
	return []byte(node.Value.(string)), nil
}
