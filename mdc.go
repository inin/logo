package logo

import (
	"fmt"
	"sync"
)

//MDC is the mapped diagnostic context for a log message
type MDC struct {
	mu   sync.RWMutex
	data map[string]string
}

//NewMDC creates a blank context
func NewMDC() *MDC {
	return &MDC{
		sync.RWMutex{},
		make(map[string]string),
	}
}

//MDCFromMDC creates a context from an existing context, copying all
//the data from the provided context.
func MDCFromMDC(ctx *MDC) *MDC {
	return &MDC{
		sync.RWMutex{},
		ctx.snapshot(),
	}
}

//MDCFromMap creates a context from the provided map. The map is not copied.
func MDCFromMap(ctx map[string]string) *MDC {
	return &MDC{
		sync.RWMutex{},
		ctx,
	}
}

//Get returns the value and existence for the given key.
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
