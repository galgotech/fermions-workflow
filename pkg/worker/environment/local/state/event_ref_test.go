package state

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
)

func Test_newEventRef(t *testing.T) {
	t.Run("run", func(t *testing.T) {
		// mapFunctions := environment.MapFunctions{
		// 	"test": &fakeFunction{},
		// }
		// test.OnEvent.EventRefs

		filterData, _ := filter.NewFilter("")
		a, err := newEventRef([]environment.Event{}, Actions{}, filterData, filterData, true)
		assert.NoError(t, err)

		dataOut, err := a.Run(context.Background(), data.Data[any]{})
		assert.Nil(t, err)
		assert.Equal(t, data.Data[any]{"test": "test"}, dataOut)
	})
}
