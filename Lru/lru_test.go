package Lru

import (
	"fmt"
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

// 测试查找函数
func TestGet(t *testing.T) {
	lru := New(20, nil)
	// 添加
	lru.Add("test1", String("test1_value"))
	if str, ok := lru.Get("test1"); ok {
		fmt.Printf("存在此内容:%s", str)
	}
}

// 测试是否会触发内存淘汰函数
func TestDelete(t *testing.T) {
	k1, k2, k3 := "test1", "test2", "test3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	// 进行添加
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Get("test1")
	lru.Add(k3, String(v3))

	str, _ := lru.Get("test1")
	t.Logf("test1 值%v", str)
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Logf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
