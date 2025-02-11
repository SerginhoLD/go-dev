package eventdispatcher

type EventDispatcher interface {
	Dispatch(event interface{})
}

type StoppableEvent interface {
	IsPropagationStopped() bool
}
