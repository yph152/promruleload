package cache

import (
	"sync"
)

type Cache struct {
	lock    *sync.RWMutex
	Rulemap map[interface{}]interface{}
}

func NewCache() *Cache {
	return &Cache{
		lock:    new(sync.RWMutex),
		Rulemap: make(map[interface{}]interface{}),
	}
}

func (m *Cache) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.Rulemap[k]; ok {
		return val
	}
	return nil
}

func (m *Cache) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.Rulemap[k]; !ok {
		m.Rulemap[k] = v
	} else if val != v {
		m.Rulemap[k] = v
	} else {
		return false
	}
	return true
}

func (m *Cache) Update(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.Rulemap[k]; !ok {
		m.Rulemap[k] = k
	} else if val != v {
		m.Rulemap[k] = v
	} else {
		return false
	}
	return true
}

func (m *Cache) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if _, ok := m.Rulemap[k]; !ok {
		return false
	}
	return true
}

func (m *Cache) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.Rulemap, k)
}
