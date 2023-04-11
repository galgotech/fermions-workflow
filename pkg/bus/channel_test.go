package bus

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	ctx := context.Background()
	channel := NewChannel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		data := <-channel.Subscribe(ctx, "test")
		assert.Equal(t, []byte("data test"), data)
		wg.Done()
	}()

	err := channel.Publish(ctx, "test", []byte("data test"))
	assert.NoError(t, err)

	wg.Wait()
}

func TestChannelWithBroadcast(t *testing.T) {
	ctx := context.Background()
	connector := NewChannel()
	broadcast := NewBroadcast(connector)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		data := <-broadcast.Subscribe(ctx, "test")
		assert.Equal(t, []byte("data test"), data)
		wg.Done()
	}()

	go func() {
		data := <-broadcast.Subscribe(ctx, "test")
		assert.Equal(t, []byte("data test"), data)
		wg.Done()
	}()

	err := connector.Publish(ctx, "test", []byte("data test"))
	assert.NoError(t, err)

	wg.Wait()
}
