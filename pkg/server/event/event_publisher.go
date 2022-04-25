package event

// EventPublish
type EventPublish struct {
	subs HashsetEventSubscribers
}

// NewEventPublish
func NewEventPublish() *EventPublish {
	pub := &EventPublish{}
	pub.subs = HashsetEventSubscribers{}

	return pub
}

// Publish
//  이밴트 리스너로 전달
func (publisher EventPublish) Publish(sender string, v ...interface{}) {
	for sub := range publisher.subs {
		sub.Update(sender, v...)
	}
}

// Subscribers
func (publisher EventPublish) Subscribers() HashsetEventSubscribers {
	return publisher.subs
}

// Close
func (publisher EventPublish) Close() {
	//clear subscribers
	for sub := range publisher.Subscribers() {
		sub.Close()
	}
}
