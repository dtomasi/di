package di

import (
	"context"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

type TestInterface interface {
	True() bool
}

type TestService1 struct {
	ctx context.Context
	c *Container
	pp ParameterProvider
	true bool
	testString string
}

func NewTestService1(ctx context.Context, c *Container,pp ParameterProvider, true bool, testString string) *TestService1 {
	return &TestService1{
		ctx: ctx,
		c: c,
		pp: pp,
		true: true,
		testString: testString,
	}
}

func (ti *TestService1) Context() context.Context {
	return ti.ctx
}

func (ti *TestService1) Container() *Container {
	return ti.c
}

func (ti *TestService1) ParamProvider() ParameterProvider {
	return ti.pp
}

func (ti *TestService1) True() bool {
	return ti.true
}

func (ti *TestService1) TestString() string {
	return ti.testString
}

type TestService2 struct {
	testService1 TestInterface
	true bool
	testString string
}

func NewTestService2(service1 TestInterface, true bool, testString string) *TestService2 {
	return &TestService2{testService1: service1, true: true, testString: testString}
}

func (ti *TestService2) True() bool {
	return ti.true
}

func (ti *TestService2) TestString() string {
	return ti.testString
}

func (ti *TestService2) TestService1() TestInterface {
	return ti.testService1
}

type ParameterProviderMock struct {}

func (m *ParameterProviderMock) Get(key string) (interface{}, error) {
	return "foo", nil
}
func (m *ParameterProviderMock) Set(key string, value interface{}) error {
	return nil
}

func BuildContainer() error {
	i := DefaultContainer()
	i.SetParameterProvider(&ParameterProviderMock{})

	i.Register(
		NewServiceDef(StringRef("TestService1")).
			Provider(NewTestService1).
			Args(
				ContextArg(),
				ContainerArg(),
				ParamProviderArg(),
				InterfaceArg(true),
				ParamArg("foo.bar.baz"),

			),

		NewServiceDef(StringRef("TestService2")).
			Provider(NewTestService2).
			Args(
				ServiceArg(StringRef("TestService1")),
				InterfaceArg(true),
				ParamArg("foo.bar.baz"),
			),
	)
	return i.Build(context.Background())
}


func TestGetContainer(t *testing.T) {
	ci := DefaultContainer()
	if ci == nil {
		t.Error("DefaultContainer returns nil value")
	}
}

func TestContainer_Build(t *testing.T) {

	err := BuildContainer()
	if err != nil {
		t.Error(err)
	}

	ci := DefaultContainer()

	t1 := ci.MustGet(StringRef("TestService1")).(*TestService1)

	assert.IsType(t, &TestService1{}, t1)
	assert.Implements(t, (*context.Context)(nil), t1.Context())
	assert.IsType(t, &Container{}, t1.Container())
	assert.Implements(t, (*ParameterProvider)(nil), t1.ParamProvider())
	assert.True(t, t1.True())
	assert.Equal(t, "foo", t1.TestString())

	t2 := ci.MustGet(StringRef("TestService2")).(*TestService2)
	assert.IsType(t, &TestService2{}, t2)
	assert.Implements(t, (*TestInterface)(nil), t2)
	assert.True(t, t2.True())
	assert.Equal(t, "foo", t2.TestString())
	assert.IsType(t, &TestService1{}, t2.TestService1())
}

func TestContainer_Build_ConcurrentRead(t *testing.T) {
	err := BuildContainer()
	if err != nil {
		t.Error(err)
	}

	for i:=0;i < 10;i++ {
		go func() {
			ci := DefaultContainer()
			for j:=0;j<10;j++ {
				_, err := ci.Get(StringRef("TestService1"))
				if err != nil {
					t.Error(err)
				}
				rand.Seed(time.Now().UnixNano())
				n := rand.Intn(100)
				time.Sleep(time.Duration(n)*time.Millisecond)
			}
		}()
	}

}