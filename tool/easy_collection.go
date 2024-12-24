package tool

import (
	"log/slog"
	"sync"
)

type ReadWriteMap[K comparable, V any] struct {
	mutex sync.RWMutex
	m     map[K]V
}

func NewReadWriteMap[K comparable, V any]() *ReadWriteMap[K, V] {
	rwMap := new(ReadWriteMap[K, V])
	rwMap.Init()
	return rwMap
}

func (rwMap *ReadWriteMap[K, V]) Init() {
	rwMap.m = make(map[K]V)
}

func (rwMap *ReadWriteMap[K, V]) Put(k K, v V) {
	rwMap.mutex.Lock()
	defer rwMap.mutex.Unlock()
	rwMap.m[k] = v
}

func (rwMap *ReadWriteMap[K, V]) Print(message string) {
	for k, v := range rwMap.m {
		slog.Info(message, "键", k, "值", v)
	}
}

func (rwMap *ReadWriteMap[K, V]) Remove(k K) {
	rwMap.mutex.Lock()
	defer rwMap.mutex.Unlock()
	delete(rwMap.m, k)
}

func (rwMap *ReadWriteMap[K, V]) Get(k K) V {
	rwMap.mutex.RLock()
	defer rwMap.mutex.RUnlock()
	return rwMap.m[k]
}
