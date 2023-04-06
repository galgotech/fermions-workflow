package bus

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadcast(t *testing.T) {
	wg := sync.WaitGroup{}
	ctx := context.Background()
	connector := &stubConnector{
		channel: make(chan []byte),
	}
	broadcast := NewBroadcast(connector)

	go func() {
		wg.Add(1)
		ch := <-broadcast.Subscribe(ctx, "test")
		assert.Equal(t, []byte("1"), ch)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		ch := <-broadcast.Subscribe(ctx, "test")
		assert.Equal(t, []byte("1"), ch)
		wg.Done()
	}()

	broadcast.Publish(ctx, "test", []byte("1"))

	wg.Wait()
}

func TestBroadcastCancel(t *testing.T) {
	wg := sync.WaitGroup{}
	ctx, ctxCancel := context.WithCancel(context.Background())
	connector := &stubConnector{
		channel: make(chan []byte),
	}
	broadcast := NewBroadcast(connector)

	go func() {
		wg.Add(1)
		ch := <-broadcast.Subscribe(ctx, "test")
		assert.Nil(t, ch)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		ch := <-broadcast.Subscribe(ctx, "test")
		assert.Nil(t, ch)
		wg.Done()
	}()

	ctxCancel()
	wg.Wait()
}
