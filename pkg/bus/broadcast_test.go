package bus

import (
	"context"
	"sync"
	"testing"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestBroadcast(t *testing.T) {
	wg := sync.WaitGroup{}
	ctx := context.Background()
	connector := &stubConnector{
		channel: make(chan []byte),
	}
	broadcast := NewBroadcast(connector)

	ch1 := broadcast.Subscribe(ctx, "test").Channel()
	ch2 := broadcast.Subscribe(ctx, "test").Channel()

	wg.Add(2)
	go func() {
		assert.Equal(t, []byte("1"), <-ch1)
		wg.Done()
	}()

	go func() {
		assert.Equal(t, []byte("1"), <-ch2)
		wg.Done()
	}()

	broadcast.Publish(ctx, "test", []byte("1"))

	wg.Wait()
}

func TestBroadcastCancelAndUnsubscribe(t *testing.T) {
	wg1 := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	wg3 := sync.WaitGroup{}
	wg4 := sync.WaitGroup{}

	ctx1, ctxCancel1 := context.WithCancel(context.Background())
	ctx2, ctxCancel2 := context.WithCancel(context.Background())
	ctx3 := context.Background()

	connector := &stubConnector{log: log.New("broacast-test"), channel: make(chan []byte)}
	broadcast := NewBroadcast(connector)
	subscribe1 := broadcast.Subscribe(ctx1, "test")
	subscribe2 := broadcast.Subscribe(ctx2, "test")
	subscribe3 := broadcast.Subscribe(ctx3, "test")

	wg1.Add(1)
	wg2.Add(1)
	wg3.Add(1)
	wg4.Add(1)

	// test context cancel
	go func() {
		assert.Nil(t, <-subscribe1.Channel())
		wg1.Done()
	}()

	// test the before context cancelation impact in this subscribe2
	go func() {
		ch := subscribe2.Channel()
		assert.Equal(t, []byte("test"), <-ch)
		wg2.Done()

		assert.Nil(t, <-ch)
		wg3.Done()
	}()

	// test unsubscribe
	go func() {
		ch := subscribe3.Channel()
		assert.Equal(t, []byte("test"), <-ch)
		assert.Nil(t, <-ch)
		wg4.Done()
	}()

	ctxCancel1()
	wg1.Wait()

	err := broadcast.Publish(ctx3, "test", []byte("test"))
	assert.NoError(t, err)
	wg2.Wait()

	ctxCancel2()
	wg3.Wait()

	subscribe3.Unsubscribe()
	wg4.Wait()
}
