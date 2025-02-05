package eventdispatcher

type EventDispatcher interface {
	Dispatch(name string, data map[string]any)
}
