package bus

import (
	"context"
	"sync"
	"time"

	"github.com/galgotech/fermions-workflow/pkg/concurrency"
	"github.com/galgotech/fermions-workflow/pkg/log"
)

func NewBroadcast(connector Connector) Connector {
	return &broadcastConnector{
		log:       log.New("test"),
		connector: connector,
		channels:  make(map[string][]chan<- []byte),
	}
}

type broadcastConnector struct {
	log       log.Logger
	connector Connector

	channelsLock sync.RWMutex
	channels     map[string][]chan<- []byte
}

func (r *broadcastConnector) Publish(ctx context.Context, name string, data []byte) error {
	return r.connector.Publish(ctx, name, data)
}

func (r *broadcastConnector) Subscribe(ctx context.Context, channelName string) <-chan []byte {
	r.broadcast(ctx, channelName)
	channel := r.addChannel(ctx, channelName)
	return channel
}

func (r *broadcastConnector) broadcast(ctx context.Context, channelName string) {
	r.channelsLock.Lock()
	defer r.channelsLock.Unlock()
	if _, ok := r.channels[channelName]; ok {
		return
	}

	srcChannel := r.connector.Subscribe(ctx, channelName)
	go func() {
		for msg := range concurrency.OrDoneCtx(ctx, srcChannel) {
			go func(msg []byte) {
				r.channelsLock.RLock()
				defer r.channelsLock.RUnlock()

				for _, ch := range r.channels[channelName] {
					go func(ch chan<- []byte) {
						select {
						case <-ctx.Done():
						case <-time.After(1 * time.Minute): // TODO: validate that timeout
						case ch <- msg:
						}
					}(ch)
				}
			}(msg)
		}
	}()
}

func (r *broadcastConnector) addChannel(ctx context.Context, channelName string) chan []byte {
	ch := make(chan []byte)

	r.channelsLock.Lock()
	r.channels[channelName] = append(r.channels[channelName], ch)
	index := len(r.channels[channelName])
	r.channelsLock.Unlock()

	go func() {
		defer close(ch)
		select {
		case <-ctx.Done():
			r.deleteChannel(channelName, index)
		}
	}()

	return ch
}

func (r *broadcastConnector) deleteChannel(channelName string, index int) {
	r.channelsLock.Lock()
	defer r.channelsLock.Unlock()
	r.channels[channelName] = r.channels[channelName][:index]
	r.channels[channelName] = append(r.channels[channelName], r.channels[channelName][index:]...)
}
