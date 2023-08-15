// utils/PriorityQueue_test.go

package utils

import (
	"container/heap"
	"testing"
	"time"
)

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue()

	// Test adding items to the queue
	items := []*Item{
		{Token: "token1", ExpiresAt: time.Now().Add(10 * time.Minute).Unix()},
		{Token: "token2", ExpiresAt: time.Now().Add(5 * time.Minute).Unix()},
		{Token: "token3", ExpiresAt: time.Now().Add(15 * time.Minute).Unix()},
	}

	for _, item := range items {
		heap.Push(pq, item) // Use heap.Push to ensure the order
	}

	// Ensure the size is as expected
	if pq.Len() != len(items) {
		t.Errorf("Expected size %d, got %d", len(items), pq.Len())
	}

	// Check if the item with the earliest expiration (token2) is at the front
	item, exists := pq.Peek()
	if !exists || item.Token != "token2" {
		t.Errorf("Expected token2 at the front, got %v", item)
	}

	// Pop the queue and ensure the order is token2, token1, token3
	expectedOrder := []string{"token2", "token1", "token3"}
	for _, expectedToken := range expectedOrder {
		item := heap.Pop(pq).(*Item) // Use heap.Pop to maintain the order
		if item.Token != expectedToken {
			t.Errorf("Expected token %s, got %s", expectedToken, item.Token)
		}
	}

	// Ensure the queue is empty
	if pq.Len() != 0 {
		t.Errorf("Expected empty queue, got size %d", pq.Len())
	}

	_, exists = pq.Peek()
	if exists {
		t.Error("Expected empty queue on peek")
	}
}

