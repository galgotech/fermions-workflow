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
	ch := channel.Subscribe(ctx, "test")
	go func() {
		assert.Equal(t, []byte("data test"), <-ch)
		wg.Done()
	}()

	err := channel.Publish(ctx, "test", []byte("data test"))
	assert.NoError(t, err)

	wg.Wait()
}

func TestChannelWithBroadcast(t *testing.T) {
	ctx := context.Background()
	broadcast := NewBroadcast(NewChannel())

	var wg sync.WaitGroup
	wg.Add(2)
	ch1 := broadcast.Subscribe(ctx, "test").Channel()
	ch2 := broadcast.Subscribe(ctx, "test").Channel()

	go func() {
		assert.Equal(t, []byte("data test"), <-ch1)
		wg.Done()
	}()
	go func() {
		assert.Equal(t, []byte("data test"), <-ch2)
		wg.Done()
	}()

	err := broadcast.Publish(ctx, "test", []byte("data test"))
	assert.NoError(t, err)

	wg.Wait()
}
