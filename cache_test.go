package main

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(2)

	// Тестируем Add и Get
	cache.Add("a", 1)
	cache.Add("b", 2)
	if val, ok := cache.Get("a"); !ok || val.(int) != 1 {
		t.Errorf("expected 1, got %v", val)
	}
	if val, ok := cache.Get("b"); !ok || val.(int) != 2 {
		t.Errorf("expected 2, got %v", val)
	}

	// Тестируем вытеснение ключей
	cache.Add("c", 3)
	if _, ok := cache.Get("a"); ok {
		t.Error("expected a to be purged")
	}

	// Тестируем Clear
	cache.Clear()
	if _, ok := cache.Get("b"); ok {
		t.Error("expected b to be cleared")
	}

	// Тестируем AddWithTTL
	cache.AddWithTTL("d", 4, 1*time.Second)
	if val, ok := cache.Get("d"); !ok || val.(int) != 4 {
		t.Errorf("expected 4, got %v", val)
	}
	time.Sleep(2 * time.Second)
	if _, ok := cache.Get("d"); ok {
		t.Error("expected d to expire")
	}

	// Тестируем Len и Cap
	cache.Add("e", 5)
	cache.Add("f", 6)
	if cache.Len() != 2 {
		t.Errorf("expected length 2, got %d", cache.Len())
	}
	if cache.Cap() != 2 {
		t.Errorf("expected capacity 2, got %d", cache.Cap())
	}

	// Тестируем Remove
	cache.Remove("e")
	if _, ok := cache.Get("e"); ok {
		t.Error("expected e to be removed")
	}
}
