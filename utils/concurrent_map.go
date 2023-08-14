// utils/concurrent_map.go

package utils

import "sync"

type ConcurrentMap struct {
	sync.RWMutex
	internal map[string]interface{}
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		internal: make(map[string]interface{}),
	}
}

func (cm *ConcurrentMap) Set(key string, value interface{}) {
	cm.Lock()
	defer cm.Unlock()
	cm.internal[key] = value
}

func (cm *ConcurrentMap) Get(key string) (interface{}, bool) {
	cm.RLock()
	defer cm.RUnlock()
	val, ok := cm.internal[key]
	return val, ok
}

func (cm *ConcurrentMap) Delete(key string) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.internal, key)
}

