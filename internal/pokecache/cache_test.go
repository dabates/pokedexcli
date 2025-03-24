package pokecache

import (
	"testing"
	"time"
)

func TestAddAndGet(t *testing.T) {
	cache := NewCache(5 * time.Second)
	key := "testkey"
	value := []byte("testvalue")

	err := cache.Add(key, value)
	if err != nil {
		t.Fatalf("expected no error from Add, got %v", err)
	}

	got, err := cache.Get(key)
	if err != nil {
		t.Fatalf("expected to find key, got error: %v", err)
	}

	if string(got) != string(value) {
		t.Errorf("expected value %s, got %s", string(value), string(got))
	}
}

func TestGetMissingKey(t *testing.T) {
	cache := NewCache(5 * time.Second)

	_, err := cache.Get("doesnotexist")
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestReapLoopRemovesExpiredEntries(t *testing.T) {
	interval := 100 * time.Millisecond
	cache := NewCache(interval)

	key := "expireme"
	value := []byte("old")

	err := cache.Add(key, value)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	time.Sleep(2 * interval) // give time for reapLoop to run

	_, err = cache.Get(key)
	if err == nil {
		t.Error("expected key to be removed by reapLoop, but it still exists")
	}
}
