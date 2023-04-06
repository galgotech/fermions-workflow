package concurrency

import (
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestOr(t *testing.T) {

	sig := func(after time.Duration) <-chan cloudevents.Event {
		c := make(chan cloudevents.Event)
		go func() {
			defer close(c)
			time.Sleep(after)
			e := cloudevents.NewEvent()

			e.SetData("duration", []byte(after.String()))
			c <- e
		}()
		return c
	}

	e := <-Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Hour),
		sig(1*time.Minute),
		sig(1*time.Second),
	)

	data := e.Data()
	assert.Equal(t, "duration", e.DataContentType())
	assert.Equal(t, "1s", string(data))
}
