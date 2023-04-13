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
	wg1 := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	wg3 := sync.WaitGroup{}
	ctx1, ctxCancel1 := context.WithCancel(context.Background())
	ctx2, ctxCancel2 := context.WithCancel(context.Background())
	ctx3 := context.Background()
	connector := &stubConnector{
		channel: make(chan []byte),
	}
	broadcast := NewBroadcast(connector)

	wg1.Add(1)
	go func() {
		ch := <-broadcast.Subscribe(ctx1, "test")
		assert.Nil(t, ch)
		wg1.Done()
	}()

	wg2.Add(1)
	wg3.Add(1)
	go func() {
		ch := broadcast.Subscribe(ctx2, "test")
		assert.Equal(t, []byte("test"), <-ch)
		wg2.Done()

		assert.Nil(t, <-ch)
		wg3.Done()
	}()

	ctxCancel1()
	wg1.Wait()

	err := broadcast.Publish(ctx3, "test", []byte("test"))
	assert.NoError(t, err)
	wg2.Wait()

	ctxCancel2()
	wg3.Wait()
}
