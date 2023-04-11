package bus

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/galgotech/fermions-workflow/pkg/concurrency"
	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/setting"
)

type Connector interface {
	Publish(ctx context.Context, name string, data []byte) error
	Subscribe(ctx context.Context, channelName string) <-chan []byte
}

type Bus interface {
	Publish(ctx context.Context, event cloudevents.Event)
	Subscribes(ctx context.Context, channelName string) <-chan BusEvent
}

func Provide(setting *setting.Setting) (Bus, error) {
	log := log.New("bus")

	bus := &BusImpl{
		log:     log,
		setting: setting,
	}

	err := bus.init()
	if err != nil {
		return nil, err
	}

	return bus, nil
}

type BusImpl struct {
	setting     *setting.Setting
	log         log.Logger
	connector   Connector
	initialized bool
}

func (b *BusImpl) init() (err error) {
	if b.initialized {
		return nil
	}
	b.initialized = true

	var connector Connector
	if b.setting.Bus.Redis != "" {
		b.log.Debug("redis url", "url", b.setting.Bus.Redis)
		connector, err = NewRedis(b.setting.Bus.Redis)
		if err != nil {
			return err
		}
	} else {
		connector = NewChannel()
	}

	b.connector = NewBroadcast(connector)
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

func (b *BusImpl) Subscribes(ctx context.Context, channel string) <-chan BusEvent {
	receive := b.connector.Subscribe(ctx, channel)
	subscribe := make(chan BusEvent)
	go func() {
		defer close(subscribe)
		for data := range concurrency.OrDoneCtx(ctx, receive) {
			busEvent := BusEvent{
				Raw: data,
			}
			event := cloudevents.NewEvent()
			err := event.UnmarshalJSON(data)
			if err != nil {
				busEvent.Err = err
			} else {
				busEvent.Event = event
			}
			subscribe <- busEvent
		}
	}()

	return subscribe
}

type BusEvent struct {
	Event cloudevents.Event
	Raw   []byte
	Err   error
}
