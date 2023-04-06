package concurrency

import "context"

func Or[T any](channels ...<-chan T) <-chan T {
	nChanels := len(channels)
	switch nChanels {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan T)
	go func() {
		defer close(orDone)
		switch nChanels {
		case 2:
			select {
			case e := <-channels[0]:
				orDone <- e
			case e := <-channels[1]:
				orDone <- e
			}
		default:
			select {
			case e := <-channels[0]:
				orDone <- e
			case e := <-channels[1]:
				orDone <- e
			case e := <-channels[2]:
				orDone <- e
			case e := <-Or(append(channels[3:], orDone)...):
				orDone <- e
			}
		}
	}()

	return orDone
}

func OrDone[T any](done <-chan bool, c <-chan T) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func OrDoneCtx[T any](ctx context.Context, c <-chan T) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
				case valStream <- v:
				}
			}
		}
	}()
	return valStream
}
