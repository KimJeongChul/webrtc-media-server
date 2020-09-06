package utils

import "sync"

// EraseSyncMap ...
func EraseSyncMap(m *sync.Map) {
	m.Range(func(key interface{}, value interface{}) bool {
		m.Delete(key)
		return true
	})
}

// EraseKeyInSyncMap ...
func EraseKeyInSyncMap(key string, m *sync.Map) {
	m.Delete(key)
}
