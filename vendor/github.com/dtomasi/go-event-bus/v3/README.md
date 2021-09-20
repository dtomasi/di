# Go EventBus

[![Go Reference](https://pkg.go.dev/badge/github.com/dtomasi/go-event-bus.svg)](https://pkg.go.dev/github.com/dtomasi/go-event-bus)
[![CodeFactor](https://www.codefactor.io/repository/github/dtomasi/go-event-bus/badge)](https://www.codefactor.io/repository/github/dtomasi/di)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/dtomasi/go-event-bus/main.svg)](https://results.pre-commit.ci/latest/github/dtomasi/go-event-bus/main)
![Go Unit Tests](https://github.com/dtomasi/go-event-bus/actions/workflows/build.yml/badge.svg)
![CodeQL](https://github.com/dtomasi/go-event-bus/actions/workflows/codeql-analysis.yml/badge.svg)
[![codecov](https://codecov.io/gh/dtomasi/go-event-bus/branch/main/graph/badge.svg?token=R8W0VJM90K)](https://codecov.io/gh/dtomasi/go-event-bus)

## Introduction

This package provides a simple yet powerful event bus.

- Simple Pub/Sub
- Async Publishing of events
- Wildcard Support

## Documentation

https://pkg.go.dev/github.com/dtomasi/go-event-bus/v3

## Installation

    go get github.com/dtomasi/go-event-bus/v3

```go
package main

import "github.com/dtomasi/go-event-bus/v3"
```

## Usage

### Simple
Subscribe and Publish events using a simple callback function

```go
package main

import "github.com/dtomasi/go-event-bus/v3"

func main()  {

    // Create a new instance
    eb := eventbus.NewEventBus()

    // Subscribe to "foo:baz" - or use a wildcard like "foo:*"
    eb.SubscribeCallback("foo:baz", func(topic string, data interface{}) {
        println(topic)
        println(data)
    })

    // Publish data to topic
    eb.Publish("foo:baz", "bar")
}
```

### Synchronous using Channels
Subscribe using a EventChannel

```go
package main

import "github.com/dtomasi/go-event-bus/v3"

func main()  {

    // Create a new instance
    eb := eventbus.NewEventBus()

    // Subscribe to "foo:baz" - or use a wildcard like "foo:*"
	eventChannel := eb.Subscribe("foo:baz")

	// Subscribe with existing channel use
	// eb.SubscribeChannel("foo:*", eventChannel)

    // Wait for the incoming event on the channel
    go func() {
        evt :=<-eventChannel
        println(evt.Topic)
        println(evt.Data)

        // Tell eventbus that you are done
        // This is only needed for synchronous publishing
        evt.Done()
    }()

    // Publish data to topic
    eb.Publish("foo:baz", "bar")
}
```

### Async
Publish asynchronously

```go
package main

import "github.com/dtomasi/go-event-bus/v3"

func main()  {

    // Create a new instance
    eb := eventbus.NewEventBus()

	// Subscribe to "foo:baz" - or use a wildcard like "foo:*"
	eventChannel := eb.Subscribe("foo:baz")

	// Subscribe with existing channel use
	// eb.SubscribeChannel("foo:*", eventChannel)

    // Wait for the incoming event on the channel
    go func() {
        evt :=<-eventChannel
        println(evt.Topic)
        println(evt.Data)
    }()

    // Publish data to topic asynchronously
    eb.PublishAsync("foo:baz", "bar")
}
```
