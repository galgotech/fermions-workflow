package bus

import (
	"context"
	"sync"
	"time"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

func NewChannel() Connector {
	return &channelConnector{
		log:      log.New("bus-channel"),
		channels: make(map[string]chan []byte),
	}
}

type channelConnector struct {
	log log.Logger

	channelsLock sync.RWMutex
	channels     map[string]chan []byte
}

func (r *channelConnector) Publish(ctx context.Context, channelName string, data []byte) error {
	channel := r.channelPubsub(channelName)
	if channel == nil {
		return nil
	}

	go func() {
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
		case channel <- data:
		}
	}()
	return nil
}

func (r *channelConnector) Subscribe(ctx context.Context, channelName string) <-chan []byte {
	ch := r.channelPubsub(channelName)
	if ch == nil {
		ch = r.createChannelPubsub(ctx, channelName)
	}
	return ch
}

func (r *channelConnector) channelPubsub(channelName string) chan []byte {
	r.channelsLock.RLock()
	defer r.channelsLock.RUnlock()
	if val, ok := r.channels[channelName]; ok {
		return val
	}
	return nil
}

func (r *channelConnector) createChannelPubsub(ctx context.Context, channelName string) chan []byte {
	r.channelsLock.Lock()
	defer r.channelsLock.Unlock()

	channel := make(chan []byte)
	r.channels[channelName] = channel

	go func() {
		select {
		case <-ctx.Done():
			r.deleteChannelPubsub(channelName)
		}
	}()

	return channel
}

func (r *channelConnector) deleteChannelPubsub(channelName string) {
	r.channelsLock.Lock()
	defer r.channelsLock.Unlock()
	close(r.channels[channelName])
	delete(r.channels, channelName)
}
