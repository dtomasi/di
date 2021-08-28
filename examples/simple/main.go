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

// SalutationParameterProvider provides us the salutations as parameter.
// This provider implements the di.ParameterProvider interface.
// As di provides this simple interface it should be relatively easy to implement providers for
// viper (github.com/spf13/viper) or koanf (github.com/knadh/koanf) or whatever you like to.
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

/*
Steps:
- create a context to provide it to di
- initialize a new di container
	- pass the context
	- add a logr interface logger for debugging
	- add our salutation prarameter provider
- register our services
- build the container
- request services and print greeting depending on daytime

example output can be found in output.txt
*/
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
		// Pass the context
		di.WithContext(ctx),
		// Add a debug logger
		di.WithLogrImpl(funcr.New(
			func(pfx, args string) { fmt.Println(pfx, args) }, //nolint:forbidigo
			funcr.Options{
				LogCaller:    funcr.All,
				LogTimestamp: true,
				Verbosity:    6, //nolint:gomnd
			}),
		),
		// Pass our parameter provider
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
