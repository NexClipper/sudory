package event

// EventPublish
type EventPublish struct {
	muxs HashsetEventNotifierMux
}

// NewEventPublish
func NewEventPublish() *EventPublish {
	pub := &EventPublish{}
	pub.muxs = HashsetEventNotifierMux{}

	return pub
}

// Publish
//  이밴트 리스너로 전달
func (publisher EventPublish) Publish(sender string, v ...interface{}) {
	for mux := range publisher.muxs {
		mux.Update(sender, v...)
	}
}

// NotifierMuxers
func (publisher EventPublish) NotifierMuxers() HashsetEventNotifierMux {
	return publisher.muxs
}

// Close
func (publisher EventPublish) Close() {
	//clear subscribers
	for mux := range publisher.NotifierMuxers() {
		mux.Close()
	}
}
