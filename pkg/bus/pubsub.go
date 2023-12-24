package bus

import (
	"context"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

type FuncNewEvent func(ctx context.Context, channelName string) Connector

func NewPubSub(newEvent FuncNewEvent) *pubSub {
	pubSub := &pubSub{
		log:         log.New("bus-pubsub"),
		publish:     make(chan newPublish),
		subscribe:   make(chan newSubscribe),
		unsubscribe: make(chan string),
		events:      make(map[string]Connector),
		newEvent:    newEvent,
	}
	pubSub.init()
	return pubSub
}

type pubSub struct {
	log         log.Logger
	publish     chan newPublish
	subscribe   chan newSubscribe
	unsubscribe chan string
	events      map[string]Connector
	newEvent    FuncNewEvent
}

func (ps *pubSub) init() {
	ps.log.Debug("pubsub initialized")
	ctx := context.Background()
	go func() {
		for {
			select {
			case <-ctx.Done():
				// TODO Clear empty channels
				ps.log.Debug("channel close")
			case newSubscribe := <-ps.subscribe:
				ps.log.Debug("subscribe", "channel", newSubscribe.channelName)
				ec, ok := ps.events[newSubscribe.channelName]
				if !ok {
					ec = ps.newEvent(ctx, newSubscribe.channelName)
					ps.events[newSubscribe.channelName] = ec
				}
				ec.Subscribe(ctx, newSubscribe.channel)
			case newPublish := <-ps.publish:
				ps.log.Debug("publish", "channel", newPublish.channelName)
				if ec, ok := ps.events[newPublish.channelName]; ok {
					ec.Publish(ctx, newPublish.data)
				} else {
					ps.log.Warn("event not defined", "channel", newPublish.channelName)
				}
			case name := <-ps.unsubscribe:
				delete(ps.events, name)
			}
		}
	}()
}

func (ps *pubSub) Subscribe(ctx context.Context, channelName string) <-chan []byte {
	channel := make(chan []byte)
	ps.subscribe <- newSubscribe{channelName, channel}
	return channel
}

func (ps *pubSub) Unsubscribe(ctx context.Context, channelName string) <-chan []byte {
	channel := make(chan []byte)
	ps.unsubscribe <- channelName
	return channel
}

func (ps *pubSub) Publish(ctx context.Context, name string, data []byte) error {
	ps.publish <- newPublish{name, data}
	return nil
}

type newSubscribe struct {
	channelName string
	channel     chan []byte
}

type newPublish struct {
	channelName string
	data        []byte
}
