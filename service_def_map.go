package di

import (
	"fmt"
	"sync"
)

type serviceMap struct {
	sync.RWMutex
	internal map[fmt.Stringer]*ServiceDef
}

func newServiceMap() *serviceMap {
	return &serviceMap{ //nolint:exhaustivestruct
		internal: map[fmt.Stringer]*ServiceDef{},
	}
}

func (rm *serviceMap) Load(key fmt.Stringer) (value *ServiceDef, ok bool) {
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

func (rm *serviceMap) Store(key fmt.Stringer, value *ServiceDef) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}

func (rm *serviceMap) Count() int {
	return len(rm.internal)
}

func (rm *serviceMap) Clear() {
	rm.Lock()
	rm.internal = map[fmt.Stringer]*ServiceDef{}
	rm.Unlock()
}

func (rm *serviceMap) Range(f func(key fmt.Stringer, def *ServiceDef) error) error {
	rm.RLock()
	for k, v := range rm.internal {
		err := f(k, v)
		if err != nil {
			return err
		}
	}
	rm.RUnlock()

	return nil
}
