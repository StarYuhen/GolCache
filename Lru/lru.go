package Lru

import (
	"container/list"
	"fmt"
)

// 2022-8-12 lru缓存策略 第一版  https://geektutu.com/post/geecache-day1.html

// LRU(Least Recently Used)
// 最近最少使用，相对于仅考虑时间因素的 FIFO 和仅考虑访问频率的 LFU，
// LRU 算法可以认为是相对平衡的一种淘汰算法。LRU 认为，如果数据最近被访问过，
// 那么将来访问的概率也会更高。LRU 算法的实现非常简单，维护一个队列，如果某条记录被访问了，
// 则移动到队尾，那么队首则是最近最少访问的数据，淘汰该条记录即可。

// list.PushFront 会将其添加到链表最后的意思时添加到了第一个添加到了最后的就是最后的内容，如添加了test1
// 如：依次添加了 "test1", "test2", "test3",写的DeleteElimination方法就会最先删除test1,因为他在最后，
// 后续的list.PushFront 都排在了第一个使用这个添加的元素后面

// Cache 定义缓存结构体
type Cache struct {
	maxByte    int64                         // 允许的最大内存
	usedByte   int64                         // 已经使用的内容
	linkedList *list.List                    // 定义的双向链表
	cache      map[string]*list.Element      // 定义的链表缓存字典
	OnRemove   func(key string, value Value) // 定义的移除回调
}

// 构建控制基础结构体,方便淘汰机制使用
type control struct {
	key   string
	value Value
}

// Value 关于值的方法
type Value interface {
	Len() int // 返回当前值的长度
}

// Len 返回链表长度--仅供test时使用
func (c *Cache) Len() int {
	return c.linkedList.Len()
}

// New 实现New创建基础缓存
func New(max int64, OnRemove func(key string, value Value)) *Cache {
	return &Cache{
		maxByte:    max,
		linkedList: list.New(),                     // 初始化双向链表
		cache:      make(map[string]*list.Element), // 使用make 返回初始化内存
		OnRemove:   OnRemove,
	}
}

// Get 查找缓存,同时将查找的双向链表节点移植最前面
func (c *Cache) Get(key string) (Value, bool) {
	// 这里利用字典的查找方法来使用，ok是因为其实Golang字典中采用查找会自动返回值和是否存在的bool
	if element, ok := c.cache[key]; ok {
		// 利用双向链表函数移植最前面
		c.linkedList.MoveToFront(element)
		kv := element.Value.(*control)
		return kv.value, true
	}
	return nil, false
}

// DeleteElimination 利用缓存淘汰来删除，即移除最近最少访问的节点（队尾）,同时更新占用缓存
func (c *Cache) DeleteElimination() {
	// back用于返回最后一个元素，也就是队尾,不存在会返回nil
	if element := c.linkedList.Back(); element != nil {
		// 删除该节点
		c.linkedList.Remove(element)
		// 删除map的映射
		kv := element.Value.(*control)
		delete(c.cache, kv.key)
		// 同时更新占用缓存
		c.usedByte -= int64(len(kv.key)) + int64(kv.value.Len())
		// 进行删除回调
		if c.OnRemove != nil {
			c.OnRemove(kv.key, kv.value)
		}
	}
}

// Add 新增或修改缓存
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		// 利用双向链表函数移植最前面
		c.linkedList.MoveToFront(element)
		kv := element.Value.(*control)
		// 更新占用缓存
		c.usedByte += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 不存在则新增,添加到双向链表最后面
		element := c.linkedList.PushFront(&control{key, value})
		c.cache[key] = element
		c.usedByte += int64(len(key)) + int64(value.Len())
	}
	fmt.Println(c.linkedList.Back().Value)
	// 当发现占用大小大于限制的最高大小则开始删除元素
	for c.maxByte != 0 && c.maxByte < c.usedByte {
		c.DeleteElimination()
	}

}
