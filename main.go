package main

import (
	"fmt"
	memCache "memory_cache/cache"
	"strconv"
	"sync"
	"time"
)

func main() {
	cache := memCache.NewMemCache()
	cache.SetMaxMemory("2kb")
	cache.Set("0", map[string]any{"1": map[string]any{"2": 3}}, 2*time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
		if _, ok := cache.Get("0"); !ok {
			fmt.Println("key(\"0\")已过期")
		}
	}()
	for i := 1; i <= 1024; i++ {
		cache.Set(strconv.Itoa(i), i, 0)
	}
	wg.Wait()
}
