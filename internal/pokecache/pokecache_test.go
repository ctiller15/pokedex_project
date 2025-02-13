package pokecache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddGetCache(t *testing.T) {
	duration := time.Duration(1 * time.Second)
	cache := NewCache(duration)

	t.Run("Successfully adds to cache", func(t *testing.T) {
		key := "testKey"
		value := []byte("testString")

		cache.Add(key, value)

		entry, found := cache.Get(key)
		if !found {
			t.Errorf("Not found")
		}

		assert.Equal(t, entry, value)
	})

	t.Run("overwrites key", func(t *testing.T) {
		key1 := "testKey1"

		val1 := []byte("testString")
		val2 := []byte("testString2")

		cache.Add(key1, val1)
		cache.Add(key1, val2)

		entry, found := cache.Get(key1)
		if !found {
			t.Errorf("Not found")
		}

		assert.Equal(t, entry, val2)
	})

	t.Run("Does not get if key is unpopulated.", func(t *testing.T) {
		key2 := "nonexistentKey"
		_, found := cache.Get(key2)
		assert.False(t, found)
	})
}

func TestReapLoop(t *testing.T) {
	t.Run("A key expires", func(t *testing.T) {
		duration := time.Duration(5 * time.Millisecond)
		cache := NewCache(duration)

		key := "testKey"
		value := []byte("testString")

		cache.Add(key, value)

		// Give it time to expire.
		time.Sleep(7 * time.Millisecond)

		_, found := cache.Get(key)
		assert.Equal(t, false, found)
	})
}
