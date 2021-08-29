package di_test

import (
	"errors"
	"fmt"
	"github.com/dtomasi/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServiceMap(t *testing.T) {
	m := di.NewServiceDefMap()
	assert.IsType(t, &di.ServiceDefMap{}, m)
}

func TestServiceDefMap_All(t *testing.T) {
	key := di.StringRef("foo")
	def := di.NewServiceDef(key)
	m := di.NewServiceDefMap()

	m.Store(key, def)
	assert.Equal(t, 1, m.Count())
	resDef, ok := m.Load(key)
	assert.True(t, ok)
	assert.Equal(t, def, resDef)

	err := m.Range(func(key fmt.Stringer, def *di.ServiceDef) error {
		assert.Implements(t, (*fmt.Stringer)(nil), key)
		assert.IsType(t, &di.ServiceDef{}, def)

		return nil
	})
	assert.NoError(t, err)

	m.Store(di.StringRef("bar"), def)
	assert.Equal(t, 2, m.Count())
	m.Delete(key)
	assert.Equal(t, 1, m.Count())
	m.Clear()
	assert.Equal(t, 0, m.Count())
}

func TestServiceDefMap_Range_Error(t *testing.T) {
	key := di.StringRef("foo")
	def := di.NewServiceDef(key)
	m := di.NewServiceDefMap()

	m.Store(key, def)

	err := m.Range(func(_ fmt.Stringer, _ *di.ServiceDef) error {
		return errors.New("something happened") // nolint:goerr113
	})

	assert.Error(t, err)
}
