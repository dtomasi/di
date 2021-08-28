package di_test

import (
	"github.com/dtomasi/di"
	"github.com/dtomasi/di/examples/simple/greeter"
)

func ExampleNewServiceDef() {
	def := di.NewServiceDef(
		// Internally di uses fmt.Stringer interface as map keys. This allows to use any type that can implement
		// this interface. You can use int´s or event pointers to a string ..
		di.StringRef("my_greeter_service_name"),
	).
		Provider(
			// Provider function that returns the service instance.
			// usually wrapped in a named function close to the service struct
			// Greeter Service borrowed from github.com/dtomasi/di/examples/simple/greeter
			greeter.NewGreeter,
		).
		Args(
			// This magic Arg injects the context provided with di container
			di.ContextArg(),
			// This magic Arg injects the logger provided with di container
			di.LoggerArg(),
			// Interface Arg can be used for any value type. The passed value must match the target value type
			// in the provider function
			di.InterfaceArg("Hello, "),
			// NOTE: There are a bunch of other magic Args available. See documentation
		).
		// Adds a tag to the service definition. This allows to group services into logical units
		AddTag(di.StringRef("greeter")).
		// Opt allows managing lifecycle of registered services within the container
		Opts(
			// This option defines that a service should not be created on Building the container.
			// It will be built when it is requested the first time via Get()
			di.BuildOnFirstRequest(),
			// This option defines that this service should always be newly created on request via Get()
			di.BuildAlwaysRebuild(),
		)

	// Getting the default (global) container instance for demo purpose
	defaultContainer := di.DefaultContainer()

	// Register the definition
	defaultContainer.Register(def)

	// Build the service container (nothing is built in our case, because of the options above)
	if err := defaultContainer.Build(); err != nil {
		panic(err)
	}

	// Get the service and greet John
	greeterService := defaultContainer.MustGet(di.StringRef("my_greeter_service_name")).(*greeter.Greeter) // nolint
	greeterService.Greet("John")
}
