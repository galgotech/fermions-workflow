package bus

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

func NewEventRedis(url string) (FuncNewEvent, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(opt)
	return func(ctx context.Context, channelName string) Connector {
		er := &eventRedis{
			log: log.New("bus-redis-event"),

			rdb:    rdb,
			pubSub: rdb.Subscribe(ctx, channelName),

			name:      channelName,
			subscribe: make(chan chan []byte),
			publish:   make(chan []byte),
			channels:  make([]chan []byte, 0),
		}
		er.init(ctx)
		return er
	}, nil
}

type eventRedis struct {
	log log.Logger

	rdb    *redis.Client
	pubSub *redis.PubSub

	name      string
	subscribe chan chan []byte
	publish   chan []byte
	channels  []chan []byte

	lock sync.RWMutex
}

func (er *eventRedis) init(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				er.log.Debug("channel close")
				for _, channel := range er.channels {
					close(channel)
				}
				er.channels = make([]chan []byte, 0)
				err := er.pubSub.Unsubscribe(ctx)
				if err != nil {
					er.log.Error("redis fail unsubscribe", "error", err.Error())
				}

			case channel := <-er.subscribe:
				er.log.Debug("event subscribe")
				er.channels = append(er.channels, channel)

			case val := <-er.pubSub.Channel():
				for i, channel := range er.channels {
					er.log.Debug("event publish", "channel-index", i)
					channel <- []byte(val.Payload)
					close(channel)
				}
				er.channels = make([]chan []byte, 0)
			}
		}
	}()
}

func (er *eventRedis) Subscribe(ctx context.Context, channel chan []byte) {
	er.subscribe <- channel
}

func (er *eventRedis) Publish(ctx context.Context, data []byte) error {
	return er.rdb.Publish(ctx, er.name, string(data)).Err()
}

func (er *eventRedis) Len() int {
	er.lock.RLock()
	defer er.lock.RUnlock()
	return len(er.channels)
}
