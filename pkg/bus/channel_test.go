package bus

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_eventChannel(t *testing.T) {
	ctx := context.Background()
	event := NewEventChannel(ctx, "test").(*eventChannel)

	for i := 0; i < 100; i++ {
		ch1 := make(chan []byte)
		ch2 := make(chan []byte)
		event.Subscribe(ctx, ch1)
		event.Subscribe(ctx, ch2)

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			assert.Equal(t, []byte("data test1"), <-ch1)
			assert.Equal(t, []uint8([]byte(nil)), <-ch1) // closed
			wg.Done()
		}()
		go func() {
			assert.Equal(t, []byte("data test1"), <-ch2)
			assert.Equal(t, []uint8([]byte(nil)), <-ch2) // closed
			wg.Done()
		}()

		event.Publish(ctx, []byte("data test1"))
		wg.Wait()

		// Empty channels
		assert.Equal(t, 0, event.Len())
		event.Publish(ctx, []byte("data test2"))
	}
}
