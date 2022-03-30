package lru

import (
	"container/list"
)

// 核心数据结构
// map + double linked list

// Cache is a LRU cache. Not safe for ocncurrent access.
type Cache struct {
	maxBytes  int64      // 允许的最大内存
	nbytes    int64      // 已使用内存
	ll        *list.List // 双向链表
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调，可以为 nil
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New 实例化
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找
//
// 如果键对应的链表结点存在，将对应节点移动到队首，并返回找到的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 缓存淘汰
//
// 移除最近最少访问的节点（尾节点）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		//TODO
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 添加一个值到缓存中
func (c *Cache) Add(key string, value Value) {
	// 已存在，修改值，并将该节点移动到队首
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len 获取已添加数据条数
func (c *Cache) Len() int {
	return c.ll.Len()
}
