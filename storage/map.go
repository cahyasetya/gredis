package storage

import (
	"fmt"
	"sync"
)

type Storage interface {
	Set(string, string)
	Get(string) (string, error)
	Del(string)
}

var (
	strg map[string]string
	mu   sync.RWMutex
)

func InitStorage() {
	if strg == nil {
		strg = make(map[string]string)
	}
}

func Set(key string, value string) {
	mu.Lock()
	defer mu.Unlock()
	strg[key] = value
}

func Get(key string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()
	val, exist := strg[key]
	if !exist {
		return "", fmt.Errorf("not found")
	}
	return val, nil
}

func Del(key string) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exist := strg[key]; !exist {
		return fmt.Errorf("key not found")
	}
	delete(strg, key)
	return nil
}


