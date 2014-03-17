package logo

import (
	"fmt"
	"sync"
)

type MDC struct {
	mu   sync.Mutex
	data map[string]string
}

func NewMDC(data map[string]string) *MDC {
	return &MDC{
		sync.Mutex{},
		data,
	}
}

func (m *MDC) Get(key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.data[key]
	return val, ok
}

func (m *MDC) Put(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = fmt.Sprintf("%v", value)
}

func (m *MDC) snapshot() map[string]string {
	m.mu.Lock()
	defer m.mu.Unlock()
	data := make(map[string]string, len(m.data))
	for key, value := range m.data {
		data[key] = value
	}
	return data
}
