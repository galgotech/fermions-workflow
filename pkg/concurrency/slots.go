package concurrency

import "context"

func Slots[T any](ctx context.Context, in <-chan T, n int) chan chan T {
	outBridge := make(chan chan T)

	go func() {
		defer close(outBridge)
		for i := 0; i < n; i++ {
			out := make(chan T)
			outBridge <- out
			go func(out chan T, i int) {
				for val := range OrDoneCtx(ctx, in) {
					out <- val
				}
			}(out, i)
		}
	}()

	return outBridge
}
