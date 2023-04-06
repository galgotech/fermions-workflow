package concurrency

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlots(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	input := make(chan int)

	outs := Slots(ctx, input, 5)

	var wg sync.WaitGroup
	var wgDone sync.WaitGroup
	for out := range outs {
		go func(out <-chan int) {
			wgDone.Add(1)
			for v := range OrDoneCtx(ctx, out) {
				assert.Equal(t, 1, v)
				wg.Done()
			}
			wgDone.Done()
		}(out)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		input <- 1
	}

	wg.Wait()
	cancel()
	wgDone.Wait()
}
