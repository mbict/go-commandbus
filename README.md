[![wercker status](https://app.wercker.com/status/f2312ebb309f222c6df110a36ca2394c/s/master "wercker status")](https://app.wercker.com/project/byKey/f2312ebb309f222c6df110a36ca2394c)
[![Build Status](https://travis-ci.org/mbict/go-commandbus.png?branch=master)](https://travis-ci.org/mbict/go-commandbus)
[![GoDoc](https://godoc.org/github.com/mbict/go-commandbus?status.png)](http://godoc.org/github.com/mbict/go-commandbus)
[![GoCover](http://gocover.io/_badge/github.com/mbict/go-commandbus)](http://gocover.io/github.com/mbict/go-commandbus)
[![GoReportCard](http://goreportcard.com/badge/mbict/go-commandbus)](http://goreportcard.com/report/mbict/go-commandbus)

# Command Bus

A simple implementation of a commandbus for go, with chain handler middleware.

#### Example
A complete example
```go
package main

import (
	"fmt"
	"context"
	"errors"
	cb "github.com/mbict/go-commandbus"
)

// DoSomethingCommand
type DoSomethingCommand struct {
	Message string
}

// CommandName is needed to identify this command
// Reflection is just too expensive for this
func (*DoSomethingCommand) CommandName() string {
	return "DoSomethingCommand"
}

func main() {
	//example of the command handler
	doSometingHandler := cb.CommandHandlerFunc(func(ctx context.Context, command cb.Command) error {
		c := command.(*DoSomethingCommand)
		fmt.Println("doSomeThingHandler says", c.Message)
		return nil
	})

	//a middleware command handler, for example validation
	middlewareHandler := cb.CommandHandlerFunc(func(_ context.Context, command cb.Command) error {
		c := command.(*DoSomethingCommand)
		if len(c.Message) < 2 {
			//if we return an error the testHandler will not be called
			return errors.New("this is no good")
		}
		return nil
	})

	//we could wrap handlers to be able to add middleware
	wrappedHandler := cb.ChainHandler(doSometingHandler, middlewareHandler)

	//creation and registration of the command
	bus := cb.New()
	bus.Register((*DoSomethingCommand)(nil), wrappedHandler)

	// create the command and call the bus
	command := &DoSomethingCommand{
		Message: "hello you!",
	}
	err := bus.Handle(nil, command)
	if err != nil {
		fmt.Println("scream outloud something failed", err)
	}
}
```