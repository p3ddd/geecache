package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	hash     Hash
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // 虚拟节点与真实节点的映射表，键是哈希值，值是真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
// 传入真实节点名称
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// fmt.Printf("key: %v, i: %v, hash: %v\n", key, i, hash)
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	// 计算 key 的哈希值
	hash := int(m.hash([]byte(key)))
	// 顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取对应的哈希值
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 得到真实节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

// func (m *Map) Remove(key string) {
// 	for i := 0; i < m.replicas; i++ {
// 		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
// 		idx := sort.SearchInts(m.keys, hash)
// 		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
// 		delete(m.hashMap, hash)
// 	}
// }
