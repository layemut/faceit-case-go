package notify

import (
	"sync"
)

type UserEvent struct {
	ID        string
	FirstName string
	LastName  string
	Email      string
}

type Pubsub struct {
	mu     sync.RWMutex
	subs   map[string][]chan UserEvent
	closed bool
}

func New() *Pubsub {
	ps := &Pubsub{}
	ps.subs = make(map[string][]chan UserEvent)
	return ps
}

func (ps *Pubsub) Subscribe(topic string) <-chan UserEvent {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan UserEvent, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *Pubsub) Publish(topic string, msg UserEvent) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}

func (ps *Pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}
