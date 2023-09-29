package events

type EventStore interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, handler SubscriptionHandler) error
}

type Event struct {
	topic string
	data  []byte
}

type SubscriptionHandler func(event Event) error

type eventStore struct {
}

func NewEventStore() {

}
