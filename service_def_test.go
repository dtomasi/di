package di

import (
	"testing"
)

type ServiceRef int

func (s ServiceRef) String() string {
	return "TestService"
}

const (
	ServiceTestRef ServiceRef = iota
)

type TestService struct {
	true bool
}

func (t *TestService) True() bool {
	return t.true
}

func TestNewServiceDef(t *testing.T) {
	sd := NewServiceDef(ServiceTestRef)
	if sd == nil {
		t.Error("NewServiceDef returns nil value")
	}
}

func TestServiceDef_build(t *testing.T) {
	sd := NewServiceDef(ServiceTestRef).
		Opts(
			BuildOnFirstRequest(),
		).
		Provider(func(foo string, true bool) (*TestService, error) {
			return &TestService{true}, nil
		}).
		Args(InterfaceArg("foo")).
		Args(InterfaceArg(true))

	// use input args as we do not need to parse hard coded values
	serviceInstance, err := sd.build(sd.args)
	if err != nil {
		t.Error(err)
	}

	serviceInstance, ok := serviceInstance.(*TestService)
	if !ok {
		t.Error("could not cast to *TestService")
	}

	wantTrue := serviceInstance.(*TestService).True()
	if !wantTrue {
		t.Error("wanted true")
	}
}
