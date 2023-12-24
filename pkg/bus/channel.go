package bus

import (
	"context"
	"sync"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

func NewEventChannel(ctx context.Context, chanelName string) Connector {
	ec := &eventChannel{
		log:       log.New("bus-event-channel"),
		name:      chanelName,
		subscribe: make(chan chan []byte),
		publish:   make(chan []byte),
		channels:  make([]chan []byte, 0),
	}
	ec.init(ctx)
	return ec
}

type eventChannel struct {
	log       log.Logger
	name      string
	subscribe chan chan []byte
	publish   chan []byte
	channels  []chan []byte
	lock      sync.RWMutex
}

func (ec *eventChannel) init(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				ec.log.Debug("channel close")
				for _, channel := range ec.channels {
					close(channel)
				}
				ec.channels = make([]chan []byte, 0)
			case channel := <-ec.subscribe:
				ec.log.Debug("event subscribe")
				ec.channels = append(ec.channels, channel)
			case publish := <-ec.publish:
				for i, channel := range ec.channels {
					ec.log.Debug("event publish", "channel-index", i)
					channel <- publish
					close(channel)
				}
				ec.channels = make([]chan []byte, 0)
			}
		}
	}()
}

func (ec *eventChannel) Subscribe(ctx context.Context, channel chan []byte) {
	ec.subscribe <- channel
}

func (ec *eventChannel) Publish(ctx context.Context, data []byte) error {
	ec.publish <- data
	return nil
}

func (ec *eventChannel) Len() int {
	ec.lock.RLock()
	defer ec.lock.RUnlock()
	return len(ec.channels)
}
