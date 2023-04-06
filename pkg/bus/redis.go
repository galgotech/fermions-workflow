package bus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/galgotech/fermions-workflow/pkg/concurrency"
	"github.com/galgotech/fermions-workflow/pkg/log"
)

func NewRedis(log log.Logger, url string) (Connector, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &redisConnector{
		log:       log,
		ctx:       context.Background(),
		subscribe: make(map[string]redisPubSup, 0),
		rdb:       redis.NewClient(opt),
	}, nil
}

var (
	Pulse   time.Duration = 20 * time.Minute
	Timeout               = 30 * time.Minute
)

type redisPubSup struct {
	pubSub *redis.PubSub
	pulse  chan<- bool
}

type redisConnector struct {
	subscribeLock sync.RWMutex
	subscribe     map[string]redisPubSup
	ctx           context.Context

	log log.Logger
	rdb *redis.Client
}

func (r *redisConnector) Publish(ctx context.Context, name string, data []byte) error {
	return r.rdb.Publish(ctx, name, string(data)).Err()
}

func (r *redisConnector) Subscribe(ctx context.Context, channelName string) <-chan []byte {
	subscribe := r.redisPubsub(channelName)
	if subscribe.pubSub == nil {
		subscribe = r.createRedisPubsub(channelName)
	}

	// pulse
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.Tick(Pulse):
				subscribe.pulse <- true
			}
		}
	}()

	channel := make(chan []byte)
	go func() {
		defer close(channel)
		for msg := range concurrency.OrDoneCtx(ctx, subscribe.pubSub.Channel()) {
			r.log.Trace("receive message", "channel", msg.Channel, "payload", msg.Payload)
			channel <- []byte(msg.Payload)
		}
	}()

	return channel
}

func (r *redisConnector) redisPubsub(channelName string) redisPubSup {
	r.subscribeLock.RLock()
	defer r.subscribeLock.RUnlock()
	if val, ok := r.subscribe[channelName]; ok {
		return val
	}
	return redisPubSup{}
}

func (r *redisConnector) createRedisPubsub(channelName string) redisPubSup {
	pubSub := r.rdb.Subscribe(r.ctx, channelName)
	pulse := make(chan bool)
	go func() {
		for {
			select {
			case <-time.After(Timeout):
				pubSub.Unsubscribe(r.ctx, channelName)
				fmt.Println("unsubscribe")
				return
			case <-r.ctx.Done():
				pubSub.Unsubscribe(r.ctx, channelName)
				return
			case <-pulse:
			}
		}
	}()

	subscribe := redisPubSup{
		pubSub: pubSub,
		pulse:  pulse,
	}

	r.subscribeLock.Lock()
	defer r.subscribeLock.Unlock()
	r.subscribe[channelName] = subscribe
	return subscribe
}

// func (r *redisConnector) redisSubscribe(ctx context.Context, channel chan<- []byte, channelsName ...string) {
// 	pubsub := r.rdb.Subscribe(ctx, channelsName...)

// 	go func() {
// 		defer pubsub.Unsubscribe(ctx)
// 		ch := pubsub.Channel()
// 		for msg := range concurrency.OrDoneCtx(ctx, ch) {
// 			r.log.Trace("receive message", "channel", msg.Channel, "payload", msg.Payload)
// 			data := []byte(msg.Payload)
// 			channel <- data
// 		}
// 	}()
// }
