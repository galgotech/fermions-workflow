package bus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	pubSub := NewPubSub(NewEventChannel)
	ctx := context.Background()

	for i := 0; i < 100; i++ {
		ch1 := pubSub.Subscribe(ctx, "test")
		ch2 := pubSub.Subscribe(ctx, "test")

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			var v []byte
			select {
			case <-time.After(1 * time.Second):
				v = []byte("")
			case v = <-ch1:
			}
			assert.Equal(t, []byte("data test"), v)
			wg.Done()
		}()
		go func() {
			var v []byte
			select {
			case <-time.After(1 * time.Second):
				v = []byte("")
			case v = <-ch2:
			}
			assert.Equal(t, []byte("data test"), v)
			wg.Done()
		}()
		pubSub.Publish(ctx, "test", []byte("data test"))
		wg.Wait()

		pubSub.Unsubscribe(ctx, "test")
	}
}
