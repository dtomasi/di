package di_test

import (
	"github.com/dtomasi/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoParameterProvider_Get(t *testing.T) {
	pp := &di.NoParameterProvider{}
	v, err := pp.Get("foo")
	assert.Nil(t, v)
	assert.NoError(t, err)
}

func TestNoParameterProvider_Set(t *testing.T) {
	pp := &di.NoParameterProvider{}
	err := pp.Set("foo", "")
	assert.NoError(t, err)
}
