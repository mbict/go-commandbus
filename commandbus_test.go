package commandbus

import (
	"context"
	"errors"
	"testing"
)

type TestCommand struct {
	Called int
}

func (_ *TestCommand) CommandName() string {
	return "test"
}

// test covers commands with the CommandName interface implemented
func TestCommandBus(t *testing.T) {
	commandHandler := CommandHandlerFunc(func(_ context.Context, command interface{}) error {
		test := command.(*TestCommand)
		test.Called++
		return nil
	})
	bus := New()
	bus.Register((*TestCommand)(nil), commandHandler)
	command := &TestCommand{}

	err := bus.Handle(nil, command)

	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	if command.Called != 1 {
		t.Fatalf("expected command to be handled once, but is called %d times", command.Called)
	}
}

type TestPlainCommand struct {
	Called int
}

//covers interface commands whos name are determined by reflection
//the command does not implement the commandName method
func TestCommandBusWithInterfaceCommand(t *testing.T) {
	commandHandler := CommandHandlerFunc(func(_ context.Context, command interface{}) error {
		test := command.(*TestPlainCommand)
		test.Called++
		return nil
	})
	bus := New()
	bus.Register((*TestPlainCommand)(nil), commandHandler)
	command := &TestPlainCommand{}

	err := bus.Handle(nil, command)

	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	if command.Called != 1 {
		t.Fatalf("expected command to be handled once, but is called %d times", command.Called)
	}
}

func TestCommandBusUnhandledCommand(t *testing.T) {
	bus := New()

	command := &TestCommand{}

	err := bus.Handle(nil, command)

	if err != ErrUnhandledCommand {
		t.Fatalf("expected error `%v`, but got %v", ErrUnhandledCommand, err)
	}
}

func TestCommandBusRegisterDuplicateHandler(t *testing.T) {
	bus := New()
	emptyHandler := CommandHandlerFunc(func(_ context.Context, command interface{}) error { return nil })

	if err := bus.Register((*TestCommand)(nil), emptyHandler); err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	err := bus.Register((*TestCommand)(nil), emptyHandler)
	if err != ErrDuplicateCommandHandler {
		t.Fatalf("expected error `%v`, but got %v", ErrDuplicateCommandHandler, err)
	}
}

type TestChainCommand struct {
	CalledByOuterHandler int
	CalledByInnerHandler int
}

func (*TestChainCommand) CommandName() string {
	return "TestChainCommand"
}

func TestCommandHandlerChain(t *testing.T) {
	innerHandler := CommandHandlerFunc(func(_ context.Context, command interface{}) error {
		test := command.(*TestChainCommand)
		test.CalledByInnerHandler++
		return nil
	})

	outerHandler := CommandHandlerFunc(func(_ context.Context, command interface{}) error {
		test := command.(*TestChainCommand)
		test.CalledByOuterHandler++
		return nil
	})

	wrappedHandler := ChainHandler(innerHandler, outerHandler)

	bus := New()
	bus.Register((*TestChainCommand)(nil), wrappedHandler)
	command := &TestChainCommand{}

	err := bus.Handle(nil, command)

	if err != nil {
		t.Fatalf("expected nil error, but got %v", err)
	}

	if command.CalledByInnerHandler != 1 {
		t.Fatalf("expected command to be handled once, but is called %d times", command.CalledByInnerHandler)
	}

	if command.CalledByOuterHandler != 1 {
		t.Fatalf("expected command to be handled once, but is called %d times", command.CalledByOuterHandler)
	}
}

func TestCommandHandlerFailEarly(t *testing.T) {
	expectedErr := errors.New("custom error")
	innerHandler := CommandHandlerFunc(func(_ context.Context, command interface{}) error {
		test := command.(*TestChainCommand)
		test.CalledByInnerHandler++
		return nil
	})

	outerHandler := CommandHandlerFunc(func(_ context.Context, command interface{}) error {
		test := command.(*TestChainCommand)
		test.CalledByOuterHandler++
		return expectedErr
	})

	wrappedHandler := ChainHandler(innerHandler, outerHandler)

	bus := New()
	bus.Register((*TestChainCommand)(nil), wrappedHandler)
	command := &TestChainCommand{}

	err := bus.Handle(nil, command)

	if err != expectedErr {
		t.Fatalf("expected error %v, but got %v", expectedErr, err)
	}

	if command.CalledByInnerHandler != 0 {
		t.Fatalf("expected command to be handled once, but is called %d times", command.CalledByInnerHandler)
	}

	if command.CalledByOuterHandler != 1 {
		t.Fatalf("expected command to be handled once, but is called %d times", command.CalledByOuterHandler)
	}
}

// DoSomethingCommand
type DoSomethingCommand struct {
	Message string
}

// CommandName is needed to identify this command
// Reflection is just too expensive for this
func (*DoSomethingCommand) CommandName() string {
	return "DoSomethingCommand"
}
