package main

import (
	"context"
	"github.com/dtomasi/di"
	"github.com/dtomasi/di/examples/simple/pkg/greeter"
)

func BuildContainer(ctx context.Context) error {
	c := di.NewServiceContainer()
	c.Register(
		// As di implements fmt.Stringer as map keys we have to pass
		// di.StringRef here to define the name (ref) of our service
		// ref could also be an int or even a pointer, as long as it
		// implements fmt.Stringer interface
		di.NewServiceDef(ServiceGoodMorningGreeter).
			// Put the reference to the provider function
			Provider(greeter.NewGreeter).
			Args(
				// This is simply a argument that allows to pass interface{}
				di.InterfaceArg("Good morning,"),
			),
		// Reuse provider to create new Good Morning Greeter
		di.NewServiceDef(ServiceGoodAfternoonGreeter).
			Opts(
				// This options tells di to not create the service on Build() but on Get()
				// We probably want to check the daytime first before greet with good morning
				di.BuildOnFirstRequest(),
			).
			Provider(greeter.NewGreeter).
			Args(
				di.InterfaceArg("Good afternoon,"),
			),

		// Reuse provider to create new Good Morning Greeter
		di.NewServiceDef(ServiceGoodEveningGreeter).
			Opts(
				// This options tells di to not create the service on Build() but on Get()
				// We probably want to check the daytime first before greet with good night
				di.BuildOnFirstRequest(),
			).
			Provider(greeter.NewGreeter).
			Args(
				di.InterfaceArg("Good evening,"),
			),
	)

	// Now that we have defined our services we want to build the container.
	// This initializes all services that are registered except the ones that we marked as "lazy"
	return c.Build(ctx)
}
