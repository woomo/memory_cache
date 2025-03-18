package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// 编译时接口实现检查
var _ Cache = (*memCache)(nil)

type memCacheValue struct {
	// 实际值
	value any
	// 过期时间，绝对时间
	expireTime time.Time
	// value 大小
	size int64
}

func (mcv *memCacheValue) isExpired() bool {
	// 非空且超时
	return !mcv.expireTime.IsZero() && time.Now().After(mcv.expireTime)
}

type memCache struct {
	// 线程锁
	mutex sync.Mutex

	// 最大内存
	maxMemorySize int64

	// 最大内存的字符串表示
	maxMemorySizeStr string

	// 当前已使用内存
	currentMemorySize int64

	// 缓存键值对的映射表
	values map[string]*memCacheValue

	// 自动触发清理过期函数的时间间隔
	clearExpiredTimeInterval time.Duration
}

func NewMemCache() Cache {
	mc := &memCache{
		mutex:                    sync.Mutex{},
		maxMemorySize:            DefaultMemSize,
		maxMemorySizeStr:         DefaultMemSizeStr,
		currentMemorySize:        0,
		clearExpiredTimeInterval: time.Second,
		values:                   make(map[string]*memCacheValue),
	}

	return mc
}

func (mc *memCache) SetMaxMemory(size string) bool {
	fmt.Println("SetMaxMemory")
	mc.maxMemorySize, mc.maxMemorySizeStr = ParseSize(size)
	fmt.Println("SetMaxMemory end", mc.maxMemorySizeStr)

	return true
}

func (mc *memCache) isLocked() bool {
	if mc.mutex.TryLock() {
		return false
	}
	return true
}

func (mc *memCache) set(key string, val *memCacheValue) {
	// 检查是否已上锁
	if !mc.isLocked() {
		return
	}
	mc.values[key] = val
	mc.currentMemorySize += val.size
}

func (mc *memCache) del(key string) {
	// 检查是否已上锁
	if !mc.isLocked() {
		return
	}
	v, ok := mc.values[key]
	if !ok {
		return
	}
	delete(mc.values, key)
	mc.currentMemorySize -= v.size
}

func (mc *memCache) cleanExpiredItems() {
	for key, val := range mc.values {
		if val.isExpired() {
			mc.del(key)
		}
	}
}

func (mc *memCache) autoCleanExpiredItems() {
	timer := time.NewTimer(mc.clearExpiredTimeInterval)
	for {
		select {
		case <-timer.C:
			mc.mutex.Lock()
			mc.cleanExpiredItems()
			mc.mutex.Unlock()
		}
	}
}

func (mc *memCache) Set(key string, val any, expiration time.Duration) bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	var previousSize int64
	if v, ok := mc.values[key]; ok {
		previousSize = v.size
	}
	size := GetValueSize(val)
	log.Printf("当前缓存大小为%d, val(%v)的大小为%d", mc.currentMemorySize, val, size)
	if mc.currentMemorySize+size-previousSize > mc.maxMemorySize {
		log.Printf("已超过最大内存缓存%s", mc.maxMemorySizeStr)
		return false
	}
	mcVal := &memCacheValue{
		value: val,
		size:  size,
	}
	if expiration > 0 {
		mcVal.expireTime = time.Now().Add(expiration)
	}
	mc.set(key, mcVal)
	return true
}

func (mc *memCache) Get(key string) (val any, exists bool) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	// 是否存在键值对
	mcv, exists := mc.values[key]
	if !exists {
		return "", false
	}
	// 值是否过期
	if mcv.isExpired() {
		mc.del(key)
		return "", false
	}
	return mcv.value, false
}

func (mc *memCache) Del(key string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.del(key)
}

func (mc *memCache) Exists(key string) bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mcv, ok := mc.values[key]
	if !ok {
		return false
	}
	if mcv.isExpired() {
		mc.del(key)
		return false
	}

	return true
}

func (mc *memCache) Flush() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.currentMemorySize = 0
	mc.values = make(map[string]*memCacheValue)
}

func (mc *memCache) Keys() int64 {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.cleanExpiredItems()

	return int64(len(mc.values))
}
