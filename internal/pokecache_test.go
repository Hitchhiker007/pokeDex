package pokecache_test

import (
	"testing"
	"time"

	pokecache "github.com/Hitchhiker007/pokeDex/internal"
)

func TestCacheAddGet(t *testing.T) {
	// cache with a short interval for testing
	cache := pokecache.NewCache(100 * time.Millisecond)

	key := "testKey"
	value := []byte("testValue")

	// add a value to the cache
	cache.Add(key, value)

	// immediately try to get the value
	got, found := cache.Get(key)
	if !found {
		t.Fatalf("expected key %s to be found", key)
	}
	if string(got) != string(value) {
		t.Errorf("expected value %s, got %s", value, got)
	}

	// wait longer than the interval to check if it gets reaped
	time.Sleep(150 * time.Millisecond)

	_, found = cache.Get(key)
	if found {
		t.Errorf("expected key %s to have expired", key)
	}
}
