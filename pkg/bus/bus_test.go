package bus

import (
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/galgotech/fermions-workflow/pkg/setting"
)

type stubSetting struct {
	setting.Setting
}

func (s *stubSetting) Bus() setting.Bus {
	return setting.Bus{
		Redis: "",
	}
}

func TestSubscribe(t *testing.T) {
	b, err := Provide(&stubSetting{})
	assert.NoError(t, err)

	ctx := context.Background()
	subscribe := b.Subscribe(ctx, "test1")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-time.After(1 * time.Second):
			assert.Fail(t, "timeout")
		case event := <-subscribe:
			assert.NoError(t, event.Err)
			d, err := event.Event.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, `{"specversion":"1.0","id":"","source":"test1","type":"type","datacontenttype":"application/json","data":{"test":"test"}}`, string(d))
		}
	}()

	event := cloudevents.NewEvent()
	event.SetType("type")
	event.SetSource("test1")
	err = event.SetData(cloudevents.ApplicationJSON, map[string]any{"test": "test"})
	assert.NoError(t, err)

	b.Publish(ctx, event)
	b.Publish(ctx, event)

	wg.Wait()
}
