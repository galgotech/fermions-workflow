package bus

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

func NewRedis(url string) (Connector, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &redisConnector{
		log:       log.New("bus-redis"),
		subscribe: make(map[string]<-chan []byte),
		rdb:       redis.NewClient(opt),
	}, nil
}

var (
	Timeout = 30 * time.Minute
)

type redisConnector struct {
	subscribeLock sync.RWMutex
	subscribe     map[string]<-chan []byte

	log log.Logger
	rdb *redis.Client
}

func (r *redisConnector) Publish(ctx context.Context, name string, data []byte) error {
	return r.rdb.Publish(ctx, name, string(data)).Err()
}

func (r *redisConnector) Subscribe(ctx context.Context, channelName string) <-chan []byte {
	channel := r.redisPubsub(channelName)
	if channel == nil {
		channel = r.createRedisPubsub(channelName)
	}

	return channel
}

func (r *redisConnector) redisPubsub(channelName string) <-chan []byte {
	r.subscribeLock.RLock()
	defer r.subscribeLock.RUnlock()
	if val, ok := r.subscribe[channelName]; ok {
		return val
	}
	return nil
}

func (r *redisConnector) createRedisPubsub(channelName string) <-chan []byte {
	ctx := context.Background()
	pubSub := r.rdb.Subscribe(ctx, channelName)

	channel := make(chan []byte)
	go func() {
		defer func() {
			close(channel)
			err := pubSub.Unsubscribe(ctx, channelName)
			if err != nil {
				r.log.Error("redis unsubscribe", "channelName", channelName, "err", err.Error())
			}
		}()

		for {
			select {
			// case timeout: TODO: add timeout and healtbeart
			case message := <-pubSub.Channel():
				channel <- []byte(message.Payload)
			}
		}
	}()

	r.subscribeLock.Lock()
	r.subscribe[channelName] = channel
	r.subscribeLock.Unlock()

	return channel
}

func (r *redisConnector) deleteRedisPubsub(channelName string) {
	r.subscribeLock.Lock()
	defer r.subscribeLock.Unlock()
	delete(r.subscribe, channelName)
}
