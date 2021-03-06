package singleflight

import "sync"

// ----------
// 一个 key 正在向其他节点、源获取数据，加锁阻塞其他相同的请求，
// 等待请求结果，防止其他节点、源压力猛增被击穿
//
// 相当于一个请求的缓存器，不具有存储功能，
// 所以请求结束后，会删除 g.m 映射关系中的 key

// 正在进行中，或已经结束的请求
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// 管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex // 保护 m 不被并发读写
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         // 如果请求正在进行中，则等待
		return c.val, c.err // 请求结束，返回结果
	}
	c := new(call)
	c.wg.Add(1)  // 发起请求前加锁
	g.m[key] = c // 添加到 g.m，表明 key 已经有对应的请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() // 调用 fn，发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key) // 更新 g.m
	g.mu.Unlock()

	return c.val, c.err // 返回结果
}
