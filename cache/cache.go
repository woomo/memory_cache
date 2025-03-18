package cache

import "time"

type Cache interface {
	// SetMaxMemory size: 1KB 100KB 1MB 2MB 1GB
	SetMaxMemory(size string) bool
	// Set 设置键值对
	Set(key string, val any, expiration time.Duration) bool
	// Get 根据key获取value
	Get(key string) (val any, exists bool)
	// Del 删除key值
	Del(key string)
	// Exists 判断key是否存在于缓存中
	Exists(key string) bool
	// Flush 清空缓存
	Flush()
	// Keys 获取缓存中key的数量
	Keys() int64
}
