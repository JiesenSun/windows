package safemap

import (
	"sync"
)

type Map struct {
	data    map[interface{}]interface{}
	rwMutex *sync.RWMutex
}

func NewMap() *Map {
	return &Map{
		data:    make(map[interface{}]interface{}),
		rwMutex: &sync.RWMutex{},
	}
}

func (m *Map) Set(key, value interface{}) {
	m.rwMutex.Lock()
	m.data[key] = value
	m.rwMutex.Unlock()
}

func (m *Map) Get(key interface{}) interface{} {
	m.rwMutex.RLock()
	v := m.data[key]
	m.rwMutex.RUnlock()
	return v
}

func (m *Map) Exist(key interface{}) bool {
	m.rwMutex.RLock()
	_, ok := m.data[key]
	m.rwMutex.RUnlock()
	return ok
}

func (m *Map) Delete(key interface{}) interface{} {
	m.rwMutex.Lock()
	v := m.data[key]
	delete(m.data, key)
	m.rwMutex.Unlock()
	return v
}
