package main

import (
	"context"
	"errors"
	"fmt"
	cb "github.com/mbict/go-commandbus/v2"
)

// DoSomethingCommand
type DoSomethingCommand struct {
	Message string
}

// CommandName is only needed to identify this command without using reflection
// It is possible to omit this method, reflection will takeover and determine the name for the command
// Reflection could be too expensive in some situations, so you can optimise by implementing the method
func (*DoSomethingCommand) CommandName() string {
	return "DoSomethingCommand"
}

func main() {
	//example of the command handler
	doSometingHandler := cb.CommandHandlerFunc(func(ctx context.Context, command interface{}) error {
		c := command.(*DoSomethingCommand)
		fmt.Println("doSomeThingHandler says", c.Message)
		return nil
	})

	//a middleware command handler, for example validation
	middlewareHandler := cb.CommandHandlerFunc(func(_ context.Context, command interface{}) error {
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
