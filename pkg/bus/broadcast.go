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
	r.broadcast(channelName)
	channel := r.addChannel(ctx, channelName)
	return channel
}

func (r *broadcastConnector) broadcast(channelName string) {
	r.channelsLock.RLock()
	defer r.channelsLock.RUnlock()
	if _, ok := r.channels[channelName]; ok {
		return
	}

	go func() {
		// TODO: Add a health to check needed keep the channelName
		ctx := context.Background()
		srcChannel := r.connector.Subscribe(ctx, channelName)
		for msg := range concurrency.OrDoneCtx(ctx, srcChannel) {
			go func(msg []byte) {
				r.channelsLock.RLock()
				defer r.channelsLock.RUnlock()

				for i, ch := range r.channels[channelName] {
					r.log.Debug("broadcast", "i", i, "channelName", channelName)
					go func(ch chan<- []byte) {
						select {
						case <-time.After(5 * time.Second): // TODO: validate that timeout
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
	r.channelsLock.Unlock()

	go func() {
		defer close(ch)
		select {
		case <-ctx.Done():
			r.deleteChannel(channelName, ch)
			r.log.Debug("broadcast delete/close channel")
		}
	}()

	return ch
}

func (r *broadcastConnector) deleteChannel(channelName string, chDelete chan<- []byte) {
	r.channelsLock.Lock()
	defer r.channelsLock.Unlock()
	index := -1
	for i, ch := range r.channels[channelName] {
		if chDelete == ch {
			index = i
			break
		}
	}

	if index == -1 {
		panic("broadcast try remove a channel not found")
	}

	r.log.Debug("remove channel from broadcast start ...", "channelName", channelName, "index", index, "len", len(r.channels[channelName]))
	start := r.channels[channelName][:index]
	end := r.channels[channelName][index+1:]
	r.channels[channelName] = append([]chan<- []byte{}, start...)
	r.channels[channelName] = append(r.channels[channelName], end...)

	r.log.Debug("remove channel from broadcast done.", "channelName", channelName, "index", index, "len", len(r.channels[channelName]))
}
