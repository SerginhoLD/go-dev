package eventdispatcher

type EventDispatcher interface {
	Dispatch(event interface{})
}
