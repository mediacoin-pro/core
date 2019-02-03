package gosync

// Trigger is asynchronous event.
type Trigger chan struct{}

// NewTrigger makes asynchronous event
func NewTrigger() Trigger {
	return make(Trigger)
}

// Done triggers the event.
// It can be called only once.
func (t Trigger) Trigger() {
	defer func() {
		recover()
	}()
	close(t)
}

func (t Trigger) Set(v bool) {
	if v {
		t.Trigger()
	}
}

// Wait waits when the event happens.
func (t Trigger) Wait() {
	<-t
}

func (t Trigger) AfterTriggeredStart(fn func()) {
	go func() {
		t.Wait()
		fn()
	}()
}
