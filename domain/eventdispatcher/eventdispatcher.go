package eventdispatcher

type EventDispatcher interface {
	Dispatch(event interface{}) error
}

type StoppableEvent interface {
	IsPropagationStopped() bool
}
