package event

type TestEvent struct {
	Value string
}

func (d *TestEvent) IsPropagationStopped() bool {
	return d.Value == "stopped" // test
}
