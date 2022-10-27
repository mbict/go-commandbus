package commandbus

import (
	"context"
	"reflect"
)

// WrapHandler is a generic wrapper for command handlers that calls the provided handler by reflection
func WrapHandler(handler interface{}) CommandHandler {
	handlerFunc := reflect.ValueOf(handler)
	return CommandHandlerFunc(func(ctx context.Context, command interface{}) error {
		results := handlerFunc.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(command),
		})

		if len(results) >= 1 {
			if results[0].IsNil() {
				return nil
			}
			return results[0].Interface().(error)
		}
		return nil
	})
}
