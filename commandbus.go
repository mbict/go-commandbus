package commandbus

import (
	"context"
	"errors"
	"reflect"
)

var ErrUnhandledCommand = errors.New("unhandled command")
var ErrDuplicateCommandHandler = errors.New("there is already a command handler registered for this command")

type Command interface {
	CommandName() string
}

type CommandHandler interface {
	Handle(ctx context.Context, command interface{}) error
}

type CommandHandlerFunc func(ctx context.Context, command interface{}) error

func (h CommandHandlerFunc) Handle(ctx context.Context, command interface{}) error {
	return h(ctx, command)
}

type CommandBus interface {
	CommandHandler
	Register(interface{}, CommandHandler) error
}

func New() CommandBus {
	return &commandBus{
		handlers: make(map[string]CommandHandler),
	}
}

// ChainHandler will wrap the handlers and call the `handler` if the `wrapper` handler succeeds
// When an error is returned by the `wrapper` handler the process will stop and return the error
func ChainHandler(handler CommandHandler, wrapper CommandHandler) CommandHandler {
	return CommandHandlerFunc(func(ctx context.Context, command interface{}) error {
		if err := wrapper.Handle(ctx, command); err != nil {
			return err
		}
		return handler.Handle(ctx, command)
	})
}

type commandBus struct {
	handlers map[string]CommandHandler
}

func (cb *commandBus) Register(command interface{}, handler CommandHandler) error {
	commandName := resolveCommandName(command)
	if _, ok := cb.handlers[commandName]; ok {
		return ErrDuplicateCommandHandler
	}
	cb.handlers[commandName] = handler
	return nil
}

func (cb *commandBus) Handle(ctx context.Context, command interface{}) error {
	if h, ok := cb.handlers[resolveCommandName(command)]; ok {
		return h.Handle(ctx, command)
	}
	return ErrUnhandledCommand
}

func resolveCommandName(command interface{}) string {
	if c, ok := command.(Command); ok {
		return c.CommandName()
	}

	t := reflect.TypeOf(command)
	if t.Kind() == reflect.Ptr {
		return t.Elem().PkgPath() + "/*" + t.Elem().Name()
	}
	return t.PkgPath() + "/" + t.Name()
}
