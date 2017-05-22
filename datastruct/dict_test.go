package datastruct

//// 测试dict go test datastruct/dict_test.go dict.go
import (
	"strconv"
	"testing"

	log "github.com/wupeaking/logrus"
)

func TestDict(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	d := NewDict(&DemoDictFuncs{})
	for i := 0; i < 100000; i++ {
		d.Set(strconv.Itoa(i), strconv.Itoa(i))
	}
	for i := 0; i < 100000; i++ {
		if v, ok := d.Get(strconv.Itoa(i)); !ok || v.(string) != strconv.Itoa(i) {
			//d.Print()
			t.Fatal("设置或者获取字典失败", "v: ", i, ok)
		}
	}

	for i := 0; i < 100000; i++ {
		d.Delete(strconv.Itoa(i))
		if v, ok := d.Get(strconv.Itoa(i)); ok {
			t.Fatal("删除字典失败:", ok, v.(string))
		}
	}

}
