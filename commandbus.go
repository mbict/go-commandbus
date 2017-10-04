package commandbus

import (
	"context"
	"errors"
)

var ErrUnhandledCommand = errors.New("unhandled command")
var ErrDuplicateCommandHandler = errors.New("there is already a command handler registered for this command")

type Command interface {
	CommandType() string
}

type CommandHandler interface {
	Handle(ctx context.Context, command Command) error
}

type CommandHandlerFunc func(ctx context.Context, command Command) error

func (h CommandHandlerFunc) Handle(ctx context.Context, command Command) error {
	return h(ctx, command)
}

type CommandBus interface {
	CommandHandler
	Register(Command, CommandHandler) error
}

func New() CommandBus {
	return &commandBus{
		handlers: make(map[string]CommandHandler),
	}
}

// ChainHandler will wrap the handlers and call the `handler` if the `wrap` handler succeeds
// When an error is returned by the `wrap` handler the process will stop and return the error
func ChainHandler(handler CommandHandler, wrap CommandHandler) CommandHandler {
	return CommandHandlerFunc(func(ctx context.Context, command Command) error {
		if err := wrap.Handle(ctx, command); err != nil {
			return err
		}
		return handler.Handle(ctx, command)
	})
}

type commandBus struct {
	handlers map[string]CommandHandler
}

func (cb *commandBus) Register(command Command, handler CommandHandler) error {
	if _, ok := cb.handlers[command.CommandType()]; ok {
		return ErrDuplicateCommandHandler
	}
	cb.handlers[command.CommandType()] = handler
	return nil
}

func (cb *commandBus) Handle(ctx context.Context, command Command) error {
	if h, ok := cb.handlers[command.CommandType()]; ok {
		return h.Handle(ctx, command)
	}
	return ErrUnhandledCommand
}
