package bus

import (
	"context"
	"errors"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/setting"
)

type Connector interface {
	Publish(ctx context.Context, data []byte) error
	Subscribe(ctx context.Context, channel chan []byte)
	Len() int
}

type BusEvent struct {
	Event cloudevents.Event
	Raw   []byte
	Err   error
}

type Bus interface {
	Publish(ctx context.Context, event cloudevents.Event)
	Subscribe(ctx context.Context, channelName string) <-chan BusEvent
}

func Provide(s setting.Setting) (Bus, error) {
	bus := &BusImpl{
		log:     log.New("bus"),
		setting: s,
	}

	err := bus.init()
	if err != nil {
		return nil, err
	}

	return bus, nil
}

type BusImpl struct {
	setting     setting.Setting
	log         log.Logger
	connector   *pubSub
	initialized bool
}

func (b *BusImpl) init() (err error) {
	if b.initialized {
		return nil
	}
	b.initialized = true

	var connector FuncNewEvent
	// TODO Add suport to https://nats.io/
	if b.setting.Bus().Redis != "" {
		b.log.Debug("redis url", "url", b.setting.Bus().Redis)
		connector, err = NewEventRedis(b.setting.Bus().Redis)
		if err != nil {
			return err
		}
	} else {
		connector = NewEventChannel
	}

	b.connector = NewPubSub(connector)
	return nil
}

func (b *BusImpl) Publish(ctx context.Context, event cloudevents.Event) {
	data, err := event.MarshalJSON()
	if err != nil {
		b.log.Error("fail marshal cloudevents", "source", event.Source())
	}

	err = b.connector.Publish(ctx, event.Source(), data)
	if err != nil {
		b.log.Error("fail publish", "event", event.Source())
	}
}

func (b *BusImpl) Subscribe(ctx context.Context, channel string) <-chan BusEvent {
	receive := b.connector.Subscribe(ctx, channel)
	subscribe := make(chan BusEvent)
	go func() {
		defer close(subscribe)
		data, ok := <-receive
		busEvent := BusEvent{
			Raw: data,
		}

		if ok {
			event := cloudevents.NewEvent()
			err := event.UnmarshalJSON(data)
			if err != nil {
				busEvent.Err = err
			} else {
				busEvent.Event = event
			}
		} else {
			busEvent.Err = errors.New("channel closed")
		}
		subscribe <- busEvent
	}()

	return subscribe
}
