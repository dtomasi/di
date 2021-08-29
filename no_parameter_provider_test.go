package di_test

import (
	"github.com/dtomasi/di"
	"testing"
)

func TestNoParameterProvider_Get(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	pp := &di.NoParameterProvider{}
	_,_ = pp.Get("foo")
}

func TestNoParameterProvider_Set(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	pp := &di.NoParameterProvider{}
	_ = pp.Set("foo", "")
}
