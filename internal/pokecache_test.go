package pokecache_test

import (
	"testing"
	"time"

	pokecache "github.com/Hitchhiker007/pokeDex/internal"
)

func TestCacheAddGet(t *testing.T) {
	// Create a cache with a short interval for testing
	cache := pokecache.NewCache(100 * time.Millisecond)

	key := "testKey"
	value := []byte("testValue")

	// Add a value to the cache
	cache.Add(key, value)

	// Immediately try to get the value
	got, found := cache.Get(key)
	if !found {
		t.Fatalf("expected key %s to be found", key)
	}
	if string(got) != string(value) {
		t.Errorf("expected value %s, got %s", value, got)
	}

	// Wait longer than the interval to ensure it gets reaped
	time.Sleep(150 * time.Millisecond)

	_, found = cache.Get(key)
	if found {
		t.Errorf("expected key %s to have expired", key)
	}
}
