package cache

import (
	"testing"
	"time"
)

func TestNewMemCache(t *testing.T) {
	cache := NewMemCache()
	if cache == nil {
		t.Fatalf("Expected new cache to be created, got nil")
	}
}

func TestSetAndGet(t *testing.T) {
	cache := NewMemCache()

	// 测试设置和获取值
	key := "key1"
	value := "value1"
	cache.Set(key, value, 0)
	val, exists := cache.Get(key)
	if !exists {
		t.Fatalf("Expected key to exist")
	}
	if val != value {
		t.Fatalf("Expected value to be %v, got %v", value, val)
	}
}

func TestExpiration(t *testing.T) {
	cache := NewMemCache()

	// 测试过期时间
	key := "key1"
	value := "value1"
	cache.Set(key, value, time.Second)

	time.Sleep(2 * time.Second)
	_, exists := cache.Get(key)
	if exists {
		t.Fatalf("Expected key to be expired")
	}
}

func TestDelete(t *testing.T) {
	cache := NewMemCache()

	// 测试删除键值对
	key := "key1"
	value := "value1"
	cache.Set(key, value, 0)
	cache.Del(key)
	_, exists := cache.Get(key)
	if exists {
		t.Fatalf("Expected key to be deleted")
	}
}

func TestExists(t *testing.T) {
	cache := NewMemCache()

	// 测试键是否存在
	key := "key1"
	value := "value1"
	cache.Set(key, value, 0)
	if !cache.Exists(key) {
		t.Fatalf("Expected key to exist")
	}

	cache.Del(key)
	if cache.Exists(key) {
		t.Fatalf("Expected key to be deleted")
	}
}

func TestFlush(t *testing.T) {
	cache := NewMemCache()

	// 测试清空缓存
	key := "key1"
	value := "value1"
	cache.Set(key, value, 0)
	cache.Flush()
	if cache.Keys() != 0 {
		t.Fatalf("Expected cache to be empty")
	}
}

func TestKeys(t *testing.T) {
	cache := NewMemCache()

	// 测试键的数量
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)
	if cache.Keys() != 2 {
		t.Fatalf("Expected 2 keys, got %d", cache.Keys())
	}
}

func TestSetMaxMemory(t *testing.T) {
	cache := NewMemCache()

	// 测试设置最大内存
	cache.SetMaxMemory("2MB")
	if cache.Set("key1", make([]byte, 3*1024*1024), 0) {
		t.Fatalf("Expected setting value to fail due to memory limit")
	}
}
