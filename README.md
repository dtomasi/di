# di (Dependency Injection)

[![CodeFactor](https://www.codefactor.io/repository/github/dtomasi/di/badge)](https://www.codefactor.io/repository/github/dtomasi/di)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/dtomasi/di/main.svg)](https://results.pre-commit.ci/latest/github/dtomasi/di/main)
![Go Unit Tests](https://github.com/dtomasi/di/actions/workflows/build.yml/badge.svg)
![CodeQL](https://github.com/dtomasi/di/actions/workflows/codeql-analysis.yml/badge.svg)


## Installation

    go get -u github.com/dtomasi/di

## Usage example

See also /examples directory

```go
package main

import (
	"context"
	"github.com/dtomasi/di"
)

func BuildContainer() error {
	i := di.DefaultContainer()

	i.Register(
		// Services are registered using fmt.Stringer interface.
		// Using this interface enables DI to use strings as well as
		// integers or even pointers as map keys.
		di.NewServiceDef(di.StringRef("TestService1")).
			// A provider function
			Provider(func() (interface{}, error) { return nil, nil}).
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
	return i.Build() //nolint:wrapcheck
}

func main() {

    err := BuildContainer()
    if err != nil {
    	panic(err)
    }

    testService := di.DefaultContainer().MustGet(di.StringRef("TestService1"))
}

```

## Licence

[Licence file](./LICENSE)
