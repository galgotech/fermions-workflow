package data

import (
	"encoding/json"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/imdario/mergo"
)

type Data[T any] map[string]T

func FromEvent(event cloudevents.Event) (Data[any], error) {
	dataOut := Data[any]{}
	err := event.DataAs(&dataOut)
	if err != nil {
		return nil, err
	}
	return dataOut, nil
}

func (d *Data[T]) Merge(src Data[T]) error {
	if err := mergo.Merge(d, src, mergo.WithOverride); err != nil {
		return err
	}
	return nil
}

func (d *Data[T]) Unmarshal(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Data[T]) FromMap(m map[string]T) {
	*d = Data[T](m)
}

func (d *Data[T]) ToMap() map[string]T {
	return map[string]T(*d)
}
