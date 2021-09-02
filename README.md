# di (Dependency Injection)

[![Go Reference](https://pkg.go.dev/badge/github.com/dtomasi/di.svg)](https://pkg.go.dev/github.com/dtomasi/di)
[![CodeFactor](https://www.codefactor.io/repository/github/dtomasi/di/badge)](https://www.codefactor.io/repository/github/dtomasi/di)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/dtomasi/di/main.svg)](https://results.pre-commit.ci/latest/github/dtomasi/di/main)
![Go Unit Tests](https://github.com/dtomasi/di/actions/workflows/build.yml/badge.svg)
![CodeQL](https://github.com/dtomasi/di/actions/workflows/codeql-analysis.yml/badge.svg)
[![codecov](https://codecov.io/gh/dtomasi/di/branch/main/graph/badge.svg?token=FBN5OAX4IK)](https://codecov.io/gh/dtomasi/di)

## Installation

    go get -u github.com/dtomasi/di

## Usage example

See also /examples directory

```go
package main

import (
	"context"
	"fmt"
	"github.com/dtomasi/di"
	"log"
)

func BuildContainer() (*di.Container, error) {
	container := di.NewServiceContainer(
		// With given context
		// di.WithContext(...)

		// With a logger that implements logr interface
		// di.WithLogrImpl(...)

		// With a parameter provider. Something like viper or koanf ...
		// di.WithParameterProvider(...)
	)

	container.Register(
		// Services are registered using fmt.Stringer interface.
		// Using this interface enables DI to use strings as well as
		// integers or even pointers as map keys.
		di.NewServiceDef(di.StringRef("TestService1")).
			// A provider function
			Provider(func() (interface{}, error) { return nil, nil }).
			// Indicated "lazy" creation of services
			Opts(di.BuildOnFirstRequest()).
			Args(
				// Injects ctx.Context from di.Build
				di.ContextArg(),
				// Injects the whole DI Container
				di.ContainerArg(),
				// Injects another service
				di.ServiceArg(di.StringRef("OtherService")),
				// Injects the registered parameter provider
				// see: parameter_provider.go
				di.ParamProviderArg(),
				// Injects a value using interface{}
				di.InterfaceArg(true),
				// Injects a parameter from registered provider
				// via Get(key string) (interface{}, error)
				di.ParamArg("foo.bar.baz"),
			),
	)

	// Builds all services
	if err := container.Build(); err != nil {
		return nil, err // nolint:wrapcheck
	}

	return container, nil
}

func main() {

	container, err := BuildContainer()
	if err != nil {
        panic(err)
	}

	testService := container.MustGet(di.StringRef("TestService1"))
}

```

## Licence

[Licence file](./LICENSE)
