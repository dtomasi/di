package di_test

import (
	"context"
	"fmt"
	"github.com/dtomasi/di"
	"github.com/go-logr/logr/funcr"
	"os"
	"os/signal"
)

func ExampleNewServiceContainer() {
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
	container := di.NewServiceContainer(
		// Pass the context
		di.WithContext(ctx),
		// Add a debug logger
		// This can be any logger that implements logr.Logger interface
		di.WithLogrImpl(funcr.New(
			func(pfx, args string) { fmt.Println(pfx, args) },
			funcr.Options{
				LogCaller:    funcr.All,
				LogTimestamp: true,
				Verbosity:    6,
			}),
		),
		// Pass our parameter provider
		di.WithParameterProvider(&di.NoParameterProvider{}),
	)

	// Build the service container
	if err := container.Build(); err != nil {
		panic(err)
	}
}
