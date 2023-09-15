package bus

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

func NewChannel(publisTimeout time.Duration) Connector {
	return &channelConnector{
		log:           log.New("bus-channel"),
		publisTimeout: publisTimeout,
		channels:      make(map[string]chan []byte),
	}
}

type channel struct {
	channelLock sync.RWMutex
	channel     map[string]chan []byte
}

type channelConnector struct {
	log           log.Logger
	publisTimeout time.Duration

	channelsLock sync.RWMutex
	channels     map[string]chan []byte
}

func (r *channelConnector) Publish(ctx context.Context, channelName string, data []byte) error {
	r.channelsLock.RLock()
	channel, ok := r.channels[channelName]
	if !ok {
		r.log.Debug("channel does not exist", "channelName", channelName)
		return nil
	}

	select {
	case <-time.After(r.publisTimeout):
		r.log.Warn("bus channel publish timeout", "channelName", channelName)
		return errors.New("bus channel publish timeout")
	case <-ctx.Done():
	case channel <- data:
	}
	r.channelsLock.RUnlock()

	return nil
}

func (r *channelConnector) Subscribe(ctx context.Context, channelName string) <-chan []byte {
	r.channelsLock.Lock()
	defer r.channelsLock.Unlock()

	channel, ok := r.channels[channelName]
	if !ok {
		channel := make(chan []byte)
		r.channels[channelName] = channel
		go func() {
			select {
			case <-ctx.Done():
				r.channelsLock.Lock()
				defer r.channelsLock.Unlock()
				ch := r.channels[channelName]
				delete(r.channels, channelName)
				close(ch)
			}
		}()
	}

	return channel
}
