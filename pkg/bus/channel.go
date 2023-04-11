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
		channel = r.createChannelPubsub(channelName)
	}

	go func() {
		select {
		case <-time.After(Timeout):
		case channel <- data:
		}

	}()
	return nil
}

func (r *channelConnector) Subscribe(ctx context.Context, channelName string) <-chan []byte {
	channel := r.channelPubsub(channelName)
	if channel == nil {
		channel = r.createChannelPubsub(channelName)
	}

	return channel
}

func (r *channelConnector) channelPubsub(channelName string) chan []byte {
	r.channelsLock.RLock()
	defer r.channelsLock.RUnlock()
	if val, ok := r.channels[channelName]; ok {
		return val
	}
	return nil
}

func (r *channelConnector) createChannelPubsub(channelName string) chan []byte {
	r.channelsLock.Lock()
	channel, ok := r.channels[channelName]
	if !ok {
		r.channels[channelName] = make(chan []byte)
	}
	channel = r.channels[channelName]
	r.channelsLock.Unlock()

	go func() {
		// defer func() {
		// 	close(channel)
		// 	r.deleteChannelPubsub(channelName)
		// }()

		// for {
		// 	select {
		// 	// case timeout: TODO: add timeout and healtbeart
		// }
	}()

	return channel
}

func (r *channelConnector) deleteChannelPubsub(channelName string) {
	r.channelsLock.Lock()
	defer r.channelsLock.Unlock()
	delete(r.channels, channelName)
}
