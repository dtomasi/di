package di

import (
	"fmt"
	"sync"
)

type serviceMap struct {
	sync.RWMutex
	internal map[fmt.Stringer]interface{}
}

func newServiceMap() *serviceMap {
	return &serviceMap{ //nolint:exhaustivestruct
		internal: map[fmt.Stringer]interface{}{},
	}
}

func (rm *serviceMap) Load(key fmt.Stringer) (value interface{}, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()

	return result, ok
}

func (rm *serviceMap) Delete(key fmt.Stringer) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *serviceMap) Store(key fmt.Stringer, value interface{}) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}

func (rm *serviceMap) Count() int {
	return len(rm.internal)
}

func (rm *serviceMap) Clear() {
	rm.Lock()
	rm.internal = map[fmt.Stringer]interface{}{}
	rm.Unlock()
}

func (rm *serviceMap) Range(f func(key, value interface{}) bool) {
	rm.RLock()
	for k, v := range rm.internal {
		retVal := f(k, v)
		if !retVal {
			break
		}
	}
	rm.RUnlock()
}
