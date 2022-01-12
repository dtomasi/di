package di_test

import (
	"context"
	"fmt"
	"github.com/dtomasi/di"
	"github.com/dtomasi/fakr"
	eventbus "github.com/dtomasi/go-event-bus/v3"
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

func BuildContainer() (*di.Container, error) {
	eb := eventbus.NewEventBus()

	eb.SubscribeCallback(di.EventTopicDIReady.String(), func(topic string, data interface{}) {
		fmt.Println("container ready") //nolint:forbidigo
	})

	container := di.NewServiceContainer(
		di.WithContext(context.Background()),
		di.WithParameterProvider(&ParameterProviderMock{}),
		di.WithLogrImpl(fakr.New()),
		di.WithEventBus(eb),
	)

	container.Register(
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
			).
			AddTag(di.StringRef("test")),
	)

	if err := container.Build(); err != nil {
		return nil, err
	}

	return container, nil
}

func TestGetContainer(t *testing.T) {
	if di.NewServiceContainer() == nil {
		t.Error("DefaultContainer returns nil value")
	}
}

func TestContainer_Build(t *testing.T) {
	container, err := BuildContainer()
	if err != nil {
		t.Error(err)
	}

	t1 := container.MustGet(di.StringRef("TestService1")).(*TestService1) //nolint:forcetypeassert

	assert.IsType(t, &TestService1{}, t1) //nolint:exhaustivestruct
	assert.Implements(t, (*context.Context)(nil), t1.Context())
	assert.IsType(t, &di.Container{}, t1.Container())
	assert.Implements(t, (*di.ParameterProvider)(nil), t1.ParamProvider())
	assert.True(t, t1.True())
	assert.Equal(t, "foo", t1.TestString())

	t2 := container.MustGet(di.StringRef("TestService2")).(*TestService2) //nolint:forcetypeassert
	assert.IsType(t, &TestService2{}, t2)                                 //nolint:exhaustivestruct
	assert.Implements(t, (*TestInterface)(nil), t2)
	assert.IsType(t, logr.Logger{}, t2.Logger())
	assert.True(t, t2.True())
	assert.Equal(t, "foo", t2.TestString())
	assert.IsType(t, &TestService1{}, t2.TestService1()) //nolint:exhaustivestruct
}

func TestContainer_Build_ConcurrentRead(t *testing.T) {
	container, err := BuildContainer()
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				_, err := container.Get(di.StringRef("TestService1"))
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

func TestContainer_Set(t *testing.T) {
	container := di.NewServiceContainer()
	container.Set(di.StringRef("foo"), &TestService1{}) // nolint:exhaustivestruct

	instance, err := container.Get(di.StringRef("foo"))
	assert.NoError(t, err)
	assert.IsType(t, &TestService1{}, instance) // nolint:exhaustivestruct
}

func TestContainer_FindByTag(t *testing.T) {
	container, err := BuildContainer()
	if err != nil {
		t.Error(err)
	}

	instances, err := container.FindByTag(di.StringRef("test"))
	assert.NoError(t, err)
	assert.Len(t, instances, 1)
	assert.IsType(t, &TestService2{}, instances[0]) // nolint:exhaustivestruct

	container.Register(di.NewServiceDef(di.StringRef("no-provider")).AddTag(di.StringRef("test")))
	_, err = container.FindByTag(di.StringRef("test"))
	assert.Error(t, err)
}

func TestContainer_GetEventBus(t *testing.T) {
	container, err := BuildContainer()
	if err != nil {
		t.Error(err)
	}

	eb := container.GetEventBus()
	assert.IsType(t, &eventbus.EventBus{}, eb)
}

func TestContainer_Get(t *testing.T) {
	container, err := BuildContainer()
	if err != nil {
		t.Error(err)
	}

	_, err = container.Get(di.StringRef("TestService1"))
	assert.NoError(t, err)

	_, err = container.Get(di.StringRef("not-exiting"))
	assert.Error(t, err)

	container.Register(di.NewServiceDef(di.StringRef("no-provider")))
	_, err = container.Get(di.StringRef("no-provider"))
	assert.Error(t, err)
}
