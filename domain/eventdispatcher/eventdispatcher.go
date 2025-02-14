package eventdispatcher

import "context"

type EventDispatcher interface {
	Dispatch(ctx context.Context, event interface{}) error
}

type StoppableEvent interface {
	IsPropagationStopped() bool
}
