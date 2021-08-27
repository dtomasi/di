package main

import (
	"context"
	"fmt"
	"github.com/dtomasi/di"
	"github.com/dtomasi/di/examples/simple/pkg/greeter"
	"github.com/goombaio/namegenerator"
	"log"
	"time"
)

var dic *di.Container

func init() {
	dic = di.DefaultContainer()
}

func main() {
	// Build the service container
	err := BuildContainer(context.Background())
	if err != nil {
		log.Fatalf("Error while building container: %v", err)
	}

	// Generate a random name
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	// Get greeter by daytime
	// NOTE: of course, this could also be achieved using a factory inside the greater service itself.
	var g *greeter.Greeter
	switch hour := time.Now().Hour(); {
	case hour < 12:
		g = getGreeterByName(ServiceGoodMorningGreeter)
	case hour < 17:
		g = getGreeterByName(ServiceGoodAfternoonGreeter)
	default:
		g = getGreeterByName(ServiceGoodEveningGreeter)
	}

	// call greet on actual greeter
	g.Greet(nameGenerator.Generate())
}

// getGreeterByName returns the requested greeter instance
func getGreeterByName(name fmt.Stringer) *greeter.Greeter {
	return dic.MustGet(name).(*greeter.Greeter)
}
