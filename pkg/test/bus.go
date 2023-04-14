package test

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/galgotech/fermions-workflow/pkg/bus"
)

func NewBusStub() bus.Bus {
	return &busStub{
		ch: make(map[string]chan bus.BusEvent, 0),
	}
}

type busStub struct {
	ch map[string]chan bus.BusEvent
}

func (b *busStub) Subscribe(ctx context.Context, channelName string) <-chan bus.BusEvent {
	ch := make(chan bus.BusEvent)
	go func() {
		channel := make(<-chan bus.BusEvent)
		if _, ok := b.ch[channelName]; !ok {
			ch := make(chan bus.BusEvent)
			b.ch[channelName] = ch
			channel = ch
		}

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-channel:
				ch <- data
			}
		}
	}()

	return ch
}

func (b *busStub) Publish(ctx context.Context, event cloudevents.Event) {
	if ch, ok := b.ch[event.Source()]; ok {
		ch <- bus.BusEvent{
			Event: event,
			Err:   nil,
		}
	}
}
