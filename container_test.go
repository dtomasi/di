package di_test

import (
	"context"
	"github.com/dtomasi/di"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

type TestInterface interface {
	True() bool
}

type TestService1 struct {
	ctx        context.Context
	c          *di.Container
	pp         di.ParameterProvider
	isTrue     bool
	testString string
}

func NewTestService1(
	ctx context.Context,
	c *di.Container,
	pp di.ParameterProvider,
	isTrue bool,
	testString string,
) *TestService1 {
	return &TestService1{
		ctx:        ctx,
		c:          c,
		pp:         pp,
		isTrue:     isTrue,
		testString: testString,
	}
}

func (ti *TestService1) Context() context.Context {
	return ti.ctx
}

func (ti *TestService1) Container() *di.Container {
	return ti.c
}

func (ti *TestService1) ParamProvider() di.ParameterProvider {
	return ti.pp
}

func (ti *TestService1) True() bool {
	return ti.isTrue
}

func (ti *TestService1) TestString() string {
	return ti.testString
}

type TestService2 struct {
	testService1 TestInterface
	logger       logr.Logger
	isTrue       bool
	testString   string
}

func NewTestService2(service1 TestInterface, logger logr.Logger, isTrue bool, testString string) *TestService2 {
	return &TestService2{testService1: service1, logger: logger, isTrue: isTrue, testString: testString}
}

func (ti *TestService2) True() bool {
	return ti.isTrue
}

func (ti *TestService2) TestString() string {
	return ti.testString
}

func (ti *TestService2) TestService1() TestInterface {
	return ti.testService1
}

func (ti *TestService2) Logger() logr.Logger {
	return ti.logger
}

type ParameterProviderMock struct{}

func (m *ParameterProviderMock) Get(_ string) (interface{}, error) {
	return "foo", nil
}
func (m *ParameterProviderMock) Set(_ string, _ interface{}) error {
	return nil
}

func BuildContainer() error {
	i := di.DefaultContainer()
	i.SetParameterProvider(&ParameterProviderMock{})

	i.Register(
		di.NewServiceDef(di.StringRef("TestService1")).
			Provider(NewTestService1).
			Args(
				di.ContextArg(),
				di.ContainerArg(),
				di.ParamProviderArg(),
				di.InterfaceArg(true),
				di.ParamArg("foo.bar.baz"),
			),

		di.NewServiceDef(di.StringRef("TestService2")).
			Provider(NewTestService2).
			Args(
				di.ServiceArg(di.StringRef("TestService1")),
				di.LoggerArg(),
				di.InterfaceArg(true),
				di.ParamArg("foo.bar.baz"),
			),
	)

	return i.Build() //nolint:wrapcheck
}

func TestGetContainer(t *testing.T) {
	if di.DefaultContainer() == nil {
		t.Error("DefaultContainer returns nil value")
	}
}

func TestContainer_Build(t *testing.T) {
	err := BuildContainer() //nolint:ifshort
	if err != nil {
		t.Error(err)
	}

	ci := di.DefaultContainer()

	t1 := ci.MustGet(di.StringRef("TestService1")).(*TestService1) //nolint:forcetypeassert

	assert.IsType(t, &TestService1{}, t1) //nolint:exhaustivestruct
	assert.Implements(t, (*context.Context)(nil), t1.Context())
	assert.IsType(t, &di.Container{}, t1.Container())
	assert.Implements(t, (*di.ParameterProvider)(nil), t1.ParamProvider())
	assert.True(t, t1.True())
	assert.Equal(t, "foo", t1.TestString())

	t2 := ci.MustGet(di.StringRef("TestService2")).(*TestService2) //nolint:forcetypeassert
	assert.IsType(t, &TestService2{}, t2)                          //nolint:exhaustivestruct
	assert.Implements(t, (*TestInterface)(nil), t2)
	assert.IsType(t, logr.Logger{}, t2.Logger())
	assert.True(t, t2.True())
	assert.Equal(t, "foo", t2.TestString())
	assert.IsType(t, &TestService1{}, t2.TestService1()) //nolint:exhaustivestruct
}

func TestContainer_Build_ConcurrentRead(t *testing.T) {
	err := BuildContainer() //nolint:ifshort
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		go func() {
			ci := di.DefaultContainer()

			for j := 0; j < 10; j++ {
				_, err := ci.Get(di.StringRef("TestService1"))
				if err != nil {
					t.Error(err)
				}

				rand.Seed(time.Now().UnixNano())
				n := rand.Intn(100) //nolint:gosec
				time.Sleep(time.Duration(n) * time.Millisecond)
			}
		}()
	}
}
