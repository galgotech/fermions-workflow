package bus

import (
	"sync"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestSubscribe(t *testing.T) {
	connector := &stubConnector{
		channel: make(chan []byte),
	}
	b := BusImpl{
		connector: connector,
	}

	event := cloudevents.NewEvent()
	event.SetType("type")
	event.SetSource("test1")
	event.UnmarshalJSON([]byte(`{}`))

	var wg sync.WaitGroup
	wg.Add(1)
	go func(t *testing.T) {
		defer wg.Done()
		ctx := context.Background()
		subscribe := b.Subscribe(ctx, "test")
		event := <-subscribe
		assert.NoError(t, event.Err)

		d, err := event.Event.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `{"specversion":"1.0","id":"","source":"test1","type":"type"}`, string(d))
	}(t)

	ctx := context.Background()
	b.Publish(ctx, event)
	b.Publish(ctx, event)
	assert.Equal(t, 3, connector.SubscribeCount)
	assert.Equal(t, 2, connector.PublishCount)

	wg.Wait()
}

type stubConnector struct {
	Connector
	PublishCount   int
	SubscribeCount int

	channel chan []byte
}

func (c *stubConnector) Publish(ctx context.Context, name string, data []byte) error {
	c.PublishCount++
	c.channel <- data
	return nil
}

func (c *stubConnector) Subscribe(ctx context.Context, channelName string) <-chan []byte {
	c.SubscribeCount = 1
	return c.channel
}
