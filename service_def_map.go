package di

import (
	"fmt"
	"sync"
)

type ServiceDefMap struct {
	mu       sync.RWMutex
	internal map[fmt.Stringer]*ServiceDef
}

func NewServiceDefMap() *ServiceDefMap {
	return &ServiceDefMap{ //nolint:exhaustivestruct
		internal: map[fmt.Stringer]*ServiceDef{},
	}
}

func (rm *ServiceDefMap) Load(key fmt.Stringer) (value *ServiceDef, ok bool) {
	rm.mu.RLock()
	result, ok := rm.internal[key]
	rm.mu.RUnlock()

	return result, ok
}

func (rm *ServiceDefMap) Delete(key fmt.Stringer) {
	rm.mu.Lock()
	delete(rm.internal, key)
	rm.mu.Unlock()
}

func (rm *ServiceDefMap) Store(key fmt.Stringer, value *ServiceDef) {
	rm.mu.Lock()
	rm.internal[key] = value
	rm.mu.Unlock()
}

func (rm *ServiceDefMap) Count() int {
	return len(rm.internal)
}

func (rm *ServiceDefMap) Clear() {
	rm.mu.Lock()
	rm.internal = map[fmt.Stringer]*ServiceDef{}
	rm.mu.Unlock()
}

func (rm *ServiceDefMap) Range(f func(key fmt.Stringer, def *ServiceDef) error) error {
	rm.mu.RLock()
	for k, v := range rm.internal {
		err := f(k, v)
		if err != nil {
			return err
		}
	}
	rm.mu.RUnlock()

	return nil
}
