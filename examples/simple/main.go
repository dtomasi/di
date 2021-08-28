package main

import (
	"context"
	"fmt"
	"github.com/dtomasi/di"
	"github.com/dtomasi/di/examples/simple/greeter"
	"github.com/go-logr/logr/funcr"
	"os"
	"os/signal"
	"time"
)

type SalutationParameterProvider struct {
	data map[string]interface{}
}

func NewMyParameterProvider() di.ParameterProvider {
	return &SalutationParameterProvider{data: map[string]interface{}{
		"morning":   "Good morning",
		"afternoon": "Good afternoon",
		"evening":   "Good evening",
	}}
}

func (p *SalutationParameterProvider) Get(key string) (interface{}, error) {
	if value, ok := p.data[key]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("key %s not found", key) // nolint
}

func (p *SalutationParameterProvider) Set(key string, value interface{}) error {
	p.data[key] = value

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	defer func() {
		signal.Stop(ch)
		cancel()
	}()

	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()

	// Create container
	c := di.NewServiceContainer(
		di.WithContext(ctx),
		di.WithLogrImpl(funcr.New(
			func(pfx, args string) { fmt.Println(pfx, args) }, //nolint:forbidigo
			funcr.Options{
				LogCaller:    funcr.All,
				LogTimestamp: true,
				Verbosity:    6, //nolint:gomnd
			}),
		),
		di.WithParameterProvider(NewMyParameterProvider()),
	)

	// register services
	c.Register(
		di.NewServiceDef(greeter.ServiceGreeterMorning).
			Opts(
				// This one is not built on calling Build()
				di.BuildOnFirstRequest(),
				).
			Provider(greeter.NewGreeter).
			Args(
				di.ParamArg("morning"),
			),
		di.NewServiceDef(greeter.ServiceGreeterAfternoon).
			Provider(greeter.NewGreeter).
			Args(
				di.ParamArg("afternoon"),
			),
		di.NewServiceDef(greeter.ServiceGreeterEvening).
			Opts(
				// This one is always recreated on request
				di.BuildAlwaysRebuild(),
			).
			Provider(greeter.NewGreeter).
			Args(
				di.ParamArg("evening"),
			),
	)

	// build container
	if err := c.Build(); err != nil {
		panic(err)
	}

	switch hour := time.Now().Hour(); {
	case hour < 12: //nolint:gomnd
		c.MustGet(greeter.ServiceGreeterMorning).(*greeter.Greeter).Greet("John")
	case hour < 17: //nolint:gomnd
		c.MustGet(greeter.ServiceGreeterAfternoon).(*greeter.Greeter).Greet("John")
	default:
		c.MustGet(greeter.ServiceGreeterEvening).(*greeter.Greeter).Greet("John")
	}
}
