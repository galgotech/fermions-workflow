package bus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	ctx := context.Background()
	channel := NewChannel(1 * time.Second)

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

func TestChannelCtxCancel(t *testing.T) {
	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		channel := NewChannel(1 * time.Second)

		channel.Subscribe(ctx, "test")
		cancel()

		err := channel.Publish(ctx, "test", []byte("data test"))
		assert.NoError(t, err)
	}

	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		channel := NewChannel(1 * time.Second)

		channel.Subscribe(ctx, "test")
		cancel()
		time.Sleep(200 * time.Millisecond)

		err := channel.Publish(ctx, "test", []byte("data test"))
		assert.NoError(t, err)
	}
}

func TestChannelWithBroadcast(t *testing.T) {
	ctx := context.Background()
	broadcast := NewBroadcast(NewChannel(1 * time.Second))

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

func TestChannelTimeout(t *testing.T) {
	ctx := context.Background()
	broadcast := NewChannel(1 * time.Millisecond)

	broadcast.Subscribe(ctx, "test")
	err := broadcast.Publish(ctx, "test", []byte("data test"))
	assert.NoError(t, err)

	time.Sleep(500 * time.Millisecond)
}
