package logo

import (
	"fmt"
	"sync"
)

type MDC struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMDC() *MDC {
	return &MDC{
		sync.RWMutex{},
		make(map[string]string),
	}
}

func MDCFromMDC(ctx *MDC) *MDC {
	return &MDC{
		sync.RWMutex{},
		ctx.snapshot(),
	}
}

func MDCFromMap(ctx map[string]string) *MDC {
	return &MDC{
		sync.RWMutex{},
		ctx,
	}
}

func (m *MDC) Get(key string) (string, bool) {
	//aquire a read lock
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.data[key]
	return val, ok
}

//Put adds the string value of 'value' to the context
func (m *MDC) Put(key string, value interface{}) {
	//aquire a write lock
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = fmt.Sprintf("%v", value)
}

func (m *MDC) snapshot() map[string]string {
	//aquire a read lock
	m.mu.RLock()
	defer m.mu.RUnlock()

	data := make(map[string]string, len(m.data))
	for key, value := range m.data {
		data[key] = value
	}
	return data
}
