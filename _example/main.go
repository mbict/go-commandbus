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

// CommandType is needed to identify this command
// Reflection is just too expensive for this
func (*DoSomethingCommand) CommandType() string {
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
	wrappedHandler := cb.ChainCommandHandler(doSometingHandler, middlewareHandler)

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
