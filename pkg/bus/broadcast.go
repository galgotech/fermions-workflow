package bus

import (
	"context"
	"sync"
	"time"

	"github.com/galgotech/fermions-workflow/pkg/concurrency"
	"github.com/galgotech/fermions-workflow/pkg/log"
)

var initialized = false

func NewBroadcast(connector Connector) *broadcastConnector {
	if initialized {
		panic("double new")
	}
	initialized = true

	broadcast := &broadcastConnector{
		log:       log.New("bus-broadcast"),
		connector: connector,
		channels:  make(map[string][]*Subscribe),
	}
	return broadcast
}

type broadcastConnector struct {
	log       log.Logger
	connector Connector

	channelsLock sync.RWMutex
	channels     map[string][]*Subscribe
}

func (r *broadcastConnector) Publish(ctx context.Context, channelName string, data []byte) error {
	r.log.Debug("publish", "channelName", channelName)
	return r.connector.Publish(ctx, channelName, data)
}

func (r *broadcastConnector) Subscribe(ctx context.Context, channelName string) *Subscribe {
	r.log.Debug("subscribe", "channelName", channelName)
	r.broadcast(ctx, channelName)
	return r.create(ctx, channelName)
}

func (r *broadcastConnector) broadcast(ctx context.Context, channelName string) {
	r.channelsLock.RLock()
	defer r.channelsLock.RUnlock()
	if _, ok := r.channels[channelName]; ok {
		return
	}

	go func() {
		ctx, ctxCancel := context.WithCancel(context.Background())
		go func() {
			// Garbage collector of source channels
			// TODO: check to add helaty check, keeping that channel
			select {
			case <-time.After(5 * time.Minute):
				r.log.Debug("broadcast gargabe collector", "channelName", channelName)

				r.channelsLock.RLock()
				l := len(r.channels[channelName])
				if l == 0 {
					ctxCancel()
				}
				r.channelsLock.RUnlock()
			}
		}()

		subscribe := r.connector.Subscribe(ctx, channelName)
		for msg := range concurrency.OrDoneCtx(ctx, subscribe) {
			r.log.Debug("source message received", "channelName", channelName)
			go func(msg []byte) {
				r.channelsLock.RLock()
				defer r.channelsLock.RUnlock()

				for i, subscribe := range r.channels[channelName] {
					r.log.Debug("broadcast to channel", "i", i, "channelName", channelName)

					go func(subscribe *Subscribe, i int) {
						// Check is closed
						if subscribe.channel == nil {
							r.log.Warn("brodcast channel closed", "index", i, "channel", channelName)
							return
						}
						select {
						case <-time.After(1 * time.Second):
							r.log.Warn("brodcast channel timeout", "index", i, "channel", channelName)
						case subscribe.channel <- msg:
						}
					}(subscribe, i)
				}
			}(msg)
		}
		r.log.Debug("broadcast finished", "channelName", channelName)
	}()
}

func (r *broadcastConnector) create(ctx context.Context, channelName string) *Subscribe {
	subscribe := &Subscribe{
		log:      r.log,
		name:     channelName,
		channel:  make(chan []byte),
		removeFn: r.remove,
	}
	go func() {
		select {
		case <-ctx.Done():
			r.log.Debug("channel context.done", "channelName", channelName)
			subscribe.Unsubscribe()
		}
	}()

	r.channelsLock.Lock()
	defer r.channelsLock.Unlock()
	r.channels[channelName] = append(r.channels[channelName], subscribe)

	return subscribe
}

func (r *broadcastConnector) remove(deleteSubscribe *Subscribe) {
	channelName := deleteSubscribe.name

	r.channelsLock.Lock()
	defer r.channelsLock.Unlock()
	for i, subscribe := range r.channels[channelName] {
		if deleteSubscribe == subscribe {
			r.log.Debug("close and remove channel from broadcast", "channelName", channelName, "index", i)

			last_index := len(r.channels[channelName]) - 1
			r.channels[channelName][i] = r.channels[channelName][last_index]
			r.channels[channelName] = r.channels[channelName][:last_index]

			close(deleteSubscribe.channel)
			break
		}
	}
}

type Subscribe struct {
	log      log.Logger
	name     string
	channel  chan []byte
	removeFn func(deleteSubscribe *Subscribe)
}

func (s *Subscribe) Channel() <-chan []byte {
	return s.channel
}

func (s *Subscribe) Unsubscribe() {
	s.log.Debug("unsubscribe", "channelName", s.channel)

	s.removeFn(s)
}
