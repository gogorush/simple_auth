// utils/concurrent_map_test.go

package utils

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentMap(t *testing.T) {
	cm := NewConcurrentMap()

	// Test Set and Get
	cm.Set("key1", "value1")
	val, exists := cm.Get("key1")
	assert.True(t, exists, "Key should exist after setting it")
	assert.Equal(t, "value1", val, "Value should match what was set")

	// Test Get for non-existent key
	_, exists = cm.Get("key2")
	assert.False(t, exists, "Key should not exist")

	// Test Delete
	cm.Delete("key1")
	_, exists = cm.Get("key1")
	assert.False(t, exists, "Key should not exist after deleting it")
}

func TestConcurrentAccess(t *testing.T) {
	cm := NewConcurrentMap()
	var wg sync.WaitGroup

	// Simulate concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cm.Set(strconv.Itoa(i), i)
		}(i)
	}

	// Simulate concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cm.Get(strconv.Itoa(i))
		}(i)
	}

	wg.Wait()

	// Ensure all values are present
	for i := 0; i < 100; i++ {
		val, exists := cm.Get(strconv.Itoa(i))
		assert.True(t, exists, "Key should exist after concurrent writes")
		assert.Equal(t, i, val, "Value should match after concurrent writes")
	}
}

